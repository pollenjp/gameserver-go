package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pollenjp/gameserver-go/api/config"
	"github.com/pollenjp/gameserver-go/api/entity"
	roomHandler "github.com/pollenjp/gameserver-go/api/handler/room"
	userHandler "github.com/pollenjp/gameserver-go/api/handler/user"
)

// - `/user/create`
func TestNewMuxUserCreate(t *testing.T) {
	t.Parallel()

	// setup
	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	mux, cleanup, err := NewMux(ctx, cfg)
	t.Cleanup(cleanup)
	if err != nil {
		t.Fatal(err)
	}

	gotBody := GotBodyOfUserCreate(t, mux, userHandler.CreateUserRequestJson{
		Name:         "test",
		LeaderCardId: 1,
	})

	CheckJsonKeyNum(t, gotBody, 1)

	// 期待する型に変換
	{
		var gotTypedJson userHandler.CreateUserResponseJson
		if err := json.Unmarshal([]byte(gotBody), &gotTypedJson); err != nil {
			t.Fatalf("json unmarshal: %v", err)
		}
	}
}

// - `/user/me`
func TestNewMuxUserMe(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	mux, cleanup, err := NewMux(ctx, cfg)
	t.Cleanup(cleanup)
	if err != nil {
		t.Fatal(err)
	}

	sampleUser := entity.User{
		Name:         "test",
		LeaderCardId: 1,
	}

	// `/user/create`
	{
		gotBody := GotBodyOfUserCreate(t, mux, userHandler.CreateUserRequestJson{
			Name:         sampleUser.Name,
			LeaderCardId: sampleUser.LeaderCardId,
		})

		var gotTypedJson userHandler.CreateUserResponseJson
		if err := json.Unmarshal([]byte(gotBody), &gotTypedJson); err != nil {
			t.Fatalf("json unmarshal: %v", err)
		}

		sampleUser.Token = gotTypedJson.Token
	}

	// `/user/me`
	{
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/user/me", nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sampleUser.Token))
		mux.ServeHTTP(w, req)
		rsp := w.Result()
		defer func() {
			_ = rsp.Body.Close()
		}()

		if rsp.StatusCode != http.StatusOK {
			t.Fatalf("status code (want %d, got %d)", http.StatusOK, rsp.StatusCode)
		}

		gotBody, err := io.ReadAll(rsp.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}

		var gotTypedJson userHandler.UserMeResponseJson
		if err := json.Unmarshal([]byte(gotBody), &gotTypedJson); err != nil {
			t.Fatalf("json unmarshal: %v", err)
		}

		expected := userHandler.UserMeResponseJson{
			Id:           gotTypedJson.Id,
			Name:         sampleUser.Name,
			LeaderCardId: sampleUser.LeaderCardId,
		}

		if diff := cmp.Diff(expected, gotTypedJson); diff != "" {
			t.Fatalf("diff: (-want +got)\n%s", diff)
		}
	}
}

// - `/user/update`
func TestNewMuxUserUpdate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	mux, cleanup, err := NewMux(ctx, cfg)
	t.Cleanup(cleanup)
	if err != nil {
		t.Fatal(err)
	}

	sampleUser := entity.User{
		Name:         "test",
		LeaderCardId: 1,
	}

	// `/user/create`
	{
		gotBody := GotBodyOfUserCreate(t, mux, userHandler.CreateUserRequestJson{
			Name:         sampleUser.Name,
			LeaderCardId: sampleUser.LeaderCardId,
		})

		var gotTypedJson userHandler.CreateUserResponseJson
		if err := json.Unmarshal([]byte(gotBody), &gotTypedJson); err != nil {
			t.Fatalf("json unmarshal: %v", err)
		}

		sampleUser.Token = gotTypedJson.Token
	}

	updatedUser := entity.User{
		Name:         "updated",
		LeaderCardId: 2,
		Token:        sampleUser.Token,
	}

	// `/user/update`
	{
		reqJsonBody := []byte(
			fmt.Sprintf(
				`{"user_name":"%s","leader_card_id":%d}`,
				updatedUser.Name,
				updatedUser.LeaderCardId,
			),
		)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/user/update", bytes.NewBuffer(reqJsonBody))
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", updatedUser.Token))
		mux.ServeHTTP(w, req)
		rsp := w.Result()
		defer func() {
			_ = rsp.Body.Close()
		}()

		gotBody, err := io.ReadAll(rsp.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}

		if rsp.StatusCode != http.StatusOK {
			FatalErrorWithStatusCodeAndBody(t, http.StatusOK, rsp.StatusCode, gotBody)
		}

		CheckJsonKeyNum(t, gotBody, 0)
	}
}

// - `/room/create`
func TestNewMuxRoomCreate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	mux, cleanup, err := NewMux(ctx, cfg)
	t.Cleanup(cleanup)
	if err != nil {
		t.Fatal(err)
	}

	sampleUser := entity.User{
		Name:         "test",
		LeaderCardId: 1,
	}

	// `/user/create`
	{
		gotBody := GotBodyOfUserCreate(t, mux, userHandler.CreateUserRequestJson{
			Name:         sampleUser.Name,
			LeaderCardId: sampleUser.LeaderCardId,
		})

		var gotTypedJson userHandler.CreateUserResponseJson
		if err := json.Unmarshal([]byte(gotBody), &gotTypedJson); err != nil {
			t.Fatalf("json unmarshal: %v", err)
		}

		sampleUser.Token = gotTypedJson.Token
	}

	// `/room/create`
	{
		reqJsonBody := []byte(
			`{
			"live_id": 1,
			"select_difficulty": 1
			}`,
		)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/room/create", bytes.NewBuffer(reqJsonBody))
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sampleUser.Token))
		mux.ServeHTTP(w, req)
		rsp := w.Result()
		defer func() {
			_ = rsp.Body.Close()
		}()

		gotBody, err := io.ReadAll(rsp.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}

		// check expected key num
		CheckJsonKeyNum(t, gotBody, 1)

		var gotTypedJson roomHandler.CreateRoomResponseJson
		if err := json.Unmarshal([]byte(gotBody), &gotTypedJson); err != nil {
			t.Fatalf("json unmarshal: %v", err)
		}

		roomId := gotTypedJson.RoomId
		if roomId <= 0 {
			t.Fatalf("room id should be greater than 0 (got %d)", roomId)
		}
	}
}

// - `/room/List`
func TestNewMuxRoomListAll(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	mux, cleanup, err := NewMux(ctx, cfg)
	t.Cleanup(cleanup)
	if err != nil {
		t.Fatal(err)
	}

	sampleUser := entity.User{
		Name:         "test",
		LeaderCardId: 1,
	}

	sampleRoom := entity.Room{
		LiveId: 1,
	}
	liveDifficulty := entity.LiveDifficultyNormal

	_, rspCreateRoom := CreateUserAndRoom(
		t,
		mux,
		userHandler.CreateUserRequestJson{
			Name:         sampleUser.Name,
			LeaderCardId: sampleUser.LeaderCardId,
		},
		roomHandler.CreateRoomRequestJson{
			LiveId:           sampleRoom.LiveId,
			SelectDifficulty: liveDifficulty,
		},
	)

	// `/room/list`

	var rspListRoom roomHandler.ListRoomResponseJson
	{
		// live_id = 0 は wildcard
		reqBody := []byte(`{
			"live_id": 0
		}`)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/room/list", bytes.NewBuffer(reqBody))
		mux.ServeHTTP(w, req)
		rsp := w.Result()
		defer func() {
			_ = rsp.Body.Close()
		}()

		gotBody, err := io.ReadAll(rsp.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}

		if rsp.StatusCode != http.StatusOK {
			FatalErrorWithStatusCodeAndBody(t, http.StatusOK, rsp.StatusCode, gotBody)
		}

		if err := json.Unmarshal([]byte(gotBody), &rspListRoom); err != nil {
			t.Fatalf("json unmarshal: %v", err)
		}
	}

	{
		roomMap := map[entity.RoomId]*roomHandler.ListRoomResponseJsonItem{}
		roomIds := []entity.RoomId{}
		for _, roomItem := range rspListRoom.RoomInfoList {
			roomIds = append(roomIds, roomItem.RoomId)
			roomMap[roomItem.RoomId] = roomItem
		}
		if _, ok := roomMap[rspCreateRoom.RoomId]; !ok {
			t.Fatalf("created room (%d) is not found in list (%v)", rspCreateRoom.RoomId, roomIds)
		}

		createdRoomItem := &roomHandler.ListRoomResponseJsonItem{
			RoomId:          rspCreateRoom.RoomId,
			LiveId:          sampleRoom.LiveId,
			JoinedUserCount: 1,
			MaxUserCount:    config.MaxUserCount,
		}
		if !reflect.DeepEqual(roomMap[rspCreateRoom.RoomId], createdRoomItem) {
			t.Fatalf("expected room item (%v), got (%v)", createdRoomItem, roomMap[rspCreateRoom.RoomId])
		}
	}
}

// - `/room/List`
func TestNewMuxRoomListFilterByLiveId(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	mux, cleanup, err := NewMux(ctx, cfg)
	t.Cleanup(cleanup)
	if err != nil {
		t.Fatal(err)
	}

	_, _ = CreateUserAndRoom(
		t,
		mux,
		userHandler.CreateUserRequestJson{
			Name:         "test user 1",
			LeaderCardId: 1,
		},
		roomHandler.CreateRoomRequestJson{
			LiveId:           entity.LiveId(1),
			SelectDifficulty: entity.LiveDifficultyNormal,
		},
	)

	_, _ = CreateUserAndRoom(
		t,
		mux,
		userHandler.CreateUserRequestJson{
			Name:         "test user 2",
			LeaderCardId: 1,
		},
		roomHandler.CreateRoomRequestJson{
			LiveId:           entity.LiveId(2),
			SelectDifficulty: entity.LiveDifficultyNormal,
		},
	)

	// `/room/list`

	filterLiveId := entity.LiveId(1)
	var rspListRoom roomHandler.ListRoomResponseJson
	{
		reqBody := []byte(fmt.Sprintf(`{
			"live_id": %d
		}`, filterLiveId))
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/room/list", bytes.NewBuffer(reqBody))
		mux.ServeHTTP(w, req)
		rsp := w.Result()
		defer func() {
			_ = rsp.Body.Close()
		}()

		gotBody, err := io.ReadAll(rsp.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}

		if rsp.StatusCode != http.StatusOK {
			FatalErrorWithStatusCodeAndBody(t, http.StatusOK, rsp.StatusCode, gotBody)
		}

		if err := json.Unmarshal([]byte(gotBody), &rspListRoom); err != nil {
			t.Fatalf("json unmarshal: %v", err)
		}
	}

	{
		roomMap := map[entity.LiveId]*roomHandler.ListRoomResponseJsonItem{}
		for _, roomItem := range rspListRoom.RoomInfoList {
			roomMap[roomItem.LiveId] = roomItem
		}

		testLiveId := entity.LiveId(2)
		if _, ok := roomMap[testLiveId]; ok {
			// testLiveId は存在しないはず
			t.Fatalf("test live id (%d) should not be found because filtered by %d", testLiveId, filterLiveId)
		}
	}
}

// - `/room/join`
func TestNewMuxRoomJoin(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		t.Fatal(err)
	}

	mux, cleanup, err := NewMux(ctx, cfg)
	t.Cleanup(cleanup)
	if err != nil {
		t.Fatal(err)
	}

	_, rspCreateRoom := CreateUserAndRoom(
		t,
		mux,
		userHandler.CreateUserRequestJson{
			Name:         "test user 1",
			LeaderCardId: 1,
		},
		roomHandler.CreateRoomRequestJson{
			LiveId:           entity.LiveId(1),
			SelectDifficulty: entity.LiveDifficultyNormal,
		},
	)

	rspCreateUserMember := CreateUser(
		t,
		mux,
		userHandler.CreateUserRequestJson{
			Name:         "test user 2",
			LeaderCardId: 1,
		},
	)

	// `/room/join`

	var rspJoinRoom roomHandler.JoinRoomResponseJson
	{
		reqBody := []byte(fmt.Sprintf(`{
			"room_id": %d,
			"select_difficulty": 1
		}`, rspCreateRoom.RoomId))
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/room/join", bytes.NewBuffer(reqBody))
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", rspCreateUserMember.Token))
		mux.ServeHTTP(w, req)
		rsp := w.Result()
		defer func() {
			_ = rsp.Body.Close()
		}()

		gotBody, err := io.ReadAll(rsp.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}

		if rsp.StatusCode != http.StatusOK {
			FatalErrorWithStatusCodeAndBody(t, http.StatusOK, rsp.StatusCode, gotBody)
		}

		if err := json.Unmarshal([]byte(gotBody), &rspJoinRoom); err != nil {
			t.Fatalf("json unmarshal: %v", err)
		}
	}

	{
		if rspJoinRoom.JoinRoomResult != entity.JoinRoomResultOk {
			t.Fatalf("expected join room result (%d), got (%d)", entity.JoinRoomResultOk, rspJoinRoom.JoinRoomResult)
		}
	}
}

// TODO: `/room/wait` (room status: waiting)
// TODO: `/room/wait` (room status: live started)
// TODO: `/room/wait` (room status: dissolution)
// TODO: `/room/start` (user: owner, room status: waiting)
// TODO: `/room/start` (user: owner, other room status)
// TODO: `/room/start` (user: not owner, other room status)
// TODO: `/room/end` (user: owner or joined user, room status: live started)
// TODO: `/room/end` (user: owner or joined user, room status: other status)
// TODO: `/room/result` (user: joined user, room status: waiting)
// TODO: `/room/result` (user: not joined user, room status: any)
// TODO: `/room/result` (最初に result を受信したあと、 n秒間待つ。n秒感を超えたら result が返ってきていないユーザーの result を0にする。)
// TODO: `/room/leave` (user: any, room status: waiting)
// TODO: `/room/leave` (user: any, room status: live started)

func FatalErrorWithStatusCodeAndBody(t *testing.T, expectedStatusCode int, gotStatusCode int, gotBody []byte) {
	t.Helper()

	t.Errorf("status code (want %d, got %d)", expectedStatusCode, gotStatusCode)
	var errorJson interface{}
	if err := json.Unmarshal([]byte(gotBody), &errorJson); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}
	t.Fatalf("error json:%v", errorJson)
}

func CheckJsonKeyNum(t *testing.T, gotBody []byte, expectedKeyNum int) {
	t.Helper()

	var gotJson interface{}
	if err := json.Unmarshal([]byte(gotBody), &gotJson); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}

	if len(gotJson.(map[string]interface{})) != expectedKeyNum {
		t.Errorf("expected to have %d key, but got %d", expectedKeyNum, len(gotJson.(map[string]interface{})))
	}
}

// response body from `/user/create`
func GotBodyOfUserCreate(t *testing.T, mux http.Handler, reqBody userHandler.CreateUserRequestJson) []byte {
	t.Helper()

	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("marshal request body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/user/create", bytes.NewBuffer(reqBodyJson))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	rsp := w.Result()
	defer func() {
		_ = rsp.Body.Close()
	}()

	gotBody, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	if rsp.StatusCode != http.StatusOK {
		FatalErrorWithStatusCodeAndBody(t, http.StatusOK, rsp.StatusCode, gotBody)
	}

	return gotBody
}

func CreateUser(
	t *testing.T,
	mux http.Handler,
	reqCreateUser userHandler.CreateUserRequestJson,
) userHandler.CreateUserResponseJson {
	t.Helper()
	gotJsonBody := GotBodyOfUserCreate(t, mux, reqCreateUser)

	var createUserResponseJson userHandler.CreateUserResponseJson
	if err := json.Unmarshal([]byte(gotJsonBody), &createUserResponseJson); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}
	return createUserResponseJson
}

func CreateUserAndRoom(
	t *testing.T,
	mux http.Handler,
	reqCreateUser userHandler.CreateUserRequestJson,
	reqCreateRoom roomHandler.CreateRoomRequestJson,
) (userHandler.CreateUserResponseJson, roomHandler.CreateRoomResponseJson) {
	t.Helper()

	// `/user/create`

	createUserResponseJson := CreateUser(t, mux, reqCreateUser)

	// `/room/create`

	var createRoomResponseJson roomHandler.CreateRoomResponseJson
	{
		reqJsonBody, err := json.Marshal(reqCreateRoom)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/room/create", bytes.NewBuffer(reqJsonBody))
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", createUserResponseJson.Token))
		mux.ServeHTTP(w, req)
		rsp := w.Result()
		defer func() {
			_ = rsp.Body.Close()
		}()

		gotBody, err := io.ReadAll(rsp.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}

		// check expected key num
		CheckJsonKeyNum(t, gotBody, 1)

		if err := json.Unmarshal([]byte(gotBody), &createRoomResponseJson); err != nil {
			t.Fatalf("json unmarshal: %v", err)
		}

		roomId := createRoomResponseJson.RoomId
		if roomId <= 0 {
			t.Fatalf("room id should be greater than 0 (got %d)", roomId)
		}
	}

	return createUserResponseJson, createRoomResponseJson
}
