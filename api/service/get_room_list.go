package service

import (
	"context"

	"github.com/pollenjp/gameserver-go/api/entity"
)

type RoomInfoItem struct {
	RoomId          entity.RoomId `db:"room_id"`
	LiveId          entity.LiveId `db:"live_id"`
	JoinedUserCount int           `db:"joined_user_count"`
}

// TODO: convert to //go:generate when writing tests
// go:generate go run github.com/matryer/moq -out get_room_list_moq_test.go . GetRoomListRepository
type GetRoomListRepository interface {
	GetRoomList(
		ctx context.Context,
		db Queryer,
		RoomStatus entity.RoomStatus,
	) ([]*RoomInfoItem, error)
	GetRoomListFilteredByLiveId(
		ctx context.Context,
		db Queryer,
		RoomStatus entity.RoomStatus,
		liveId entity.LiveId,
	) ([]*RoomInfoItem, error)
}

type GetRoomList struct {
	DB   Queryer
	Repo GetRoomListRepository
}

func (ru *GetRoomList) GetRoomList(
	ctx context.Context,
	liveId entity.LiveId,
) ([]*RoomInfoItem, error) {
	var roomItemList []*RoomInfoItem
	{
		var err error
		if liveId == entity.LiveId(0) {
			roomItemList, err = ru.Repo.GetRoomList(ctx, ru.DB, entity.RoomStatusWaiting)
		} else {
			roomItemList, err = ru.Repo.GetRoomListFilteredByLiveId(ctx, ru.DB, entity.RoomStatusWaiting, liveId)
		}
		if err != nil {
			return nil, err
		}
	}
	return roomItemList, nil
}
