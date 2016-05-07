// Code generated by protoc-gen-go.
// source: contextManager.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	contextManager.proto

It has these top-level messages:
	Team
	Teams
	Player
	Players
	Game
	Games
	BatchUpdates
	Result
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto1.ProtoPackageIsVersion1

type Result_Status int32

const (
	Result_SUCCESS Result_Status = 0
	Result_FAIL    Result_Status = 1
)

var Result_Status_name = map[int32]string{
	0: "SUCCESS",
	1: "FAIL",
}
var Result_Status_value = map[string]int32{
	"SUCCESS": 0,
	"FAIL":    1,
}

func (x Result_Status) String() string {
	return proto1.EnumName(Result_Status_name, int32(x))
}
func (Result_Status) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{7, 0} }

type Team struct {
	Teamid int32  `protobuf:"varint,1,opt,name=teamid" json:"teamid,omitempty"`
	Lcsid  int32  `protobuf:"varint,2,opt,name=lcsid" json:"lcsid,omitempty"`
	Riotid int32  `protobuf:"varint,3,opt,name=riotid" json:"riotid,omitempty"`
	Name   string `protobuf:"bytes,4,opt,name=name" json:"name,omitempty"`
	Tag    string `protobuf:"bytes,5,opt,name=tag" json:"tag,omitempty"`
}

fteamunc (m *Team) Reset()                    { *m = Team{} }
func (m *Team) String() string            { return proto1.CompactTextString(m) }
func (*Team) ProtoMessage()               {}
func (*Team) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Teams struct {
	Teams []*Team `protobuf:"bytes,1,rep,name=teams" json:"teams,omitempty"`
}

func (m *Teams) Reset()                    { *m = Teams{} }
func (m *Teams) String() string            { return proto1.CompactTextString(m) }
func (*Teams) ProtoMessage()               {}
func (*Teams) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Teams) GetTeams() []*Team {
	if m != nil {
		return m.Teams
	}
	return nil
}

type Player struct {
	Playerid int32  `protobuf:"varint,1,opt,name=playerid" json:"playerid,omitempty"`
	Lcsid    int32  `protobuf:"varint,2,opt,name=lcsid" json:"lcsid,omitempty"`
	Riotid   int32  `protobuf:"varint,3,opt,name=riotid" json:"riotid,omitempty"`
	Name     string `protobuf:"bytes,4,opt,name=name" json:"name,omitempty"`
	Teamid   int32  `protobuf:"varint,5,opt,name=teamid" json:"teamid,omitempty"`
}

func (m *Player) Reset()                    { *m = Player{} }
func (m *Player) String() string            { return proto1.CompactTextString(m) }
func (*Player) ProtoMessage()               {}
func (*Player) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type Players struct {
	Players []*Player `protobuf:"bytes,1,rep,name=players" json:"players,omitempty"`
}

func (m *Players) Reset()                    { *m = Players{} }
func (m *Players) String() string            { return proto1.CompactTextString(m) }
func (*Players) ProtoMessage()               {}
func (*Players) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Players) GetPlayers() []*Player {
	if m != nil {
		return m.Players
	}
	return nil
}

type Game struct {
	Gameid      int32  `protobuf:"varint,1,opt,name=gameid" json:"gameid,omitempty"`
	Lcsid       int32  `protobuf:"varint,2,opt,name=lcsid" json:"lcsid,omitempty"`
	Riotgameid  string `protobuf:"bytes,3,opt,name=riotgameid" json:"riotgameid,omitempty"`
	Riotmatchid string `protobuf:"bytes,4,opt,name=riotmatchid" json:"riotmatchid,omitempty"`
	Redteamid   int32  `protobuf:"varint,5,opt,name=redteamid" json:"redteamid,omitempty"`
	Blueteamid  int32  `protobuf:"varint,6,opt,name=blueteamid" json:"blueteamid,omitempty"`
	Gamestart   int64  `protobuf:"varint,7,opt,name=gamestart" json:"gamestart,omitempty"`
	Gameend     int64  `protobuf:"varint,8,opt,name=gameend" json:"gameend,omitempty"`
}

func (m *Game) Reset()                    { *m = Game{} }
func (m *Game) String() string            { return proto1.CompactTextString(m) }
func (*Game) ProtoMessage()               {}
func (*Game) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type Games struct {
	Games []*Game `protobuf:"bytes,1,rep,name=games" json:"games,omitempty"`
}

func (m *Games) Reset()                    { *m = Games{} }
func (m *Games) String() string            { return proto1.CompactTextString(m) }
func (*Games) ProtoMessage()               {}
func (*Games) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Games) GetGames() []*Game {
	if m != nil {
		return m.Games
	}
	return nil
}

type BatchUpdates struct {
	TeamsCreate   *Teams   `protobuf:"bytes,1,opt,name=teamsCreate" json:"teamsCreate,omitempty"`
	TeamsUpdate   *Teams   `protobuf:"bytes,2,opt,name=teamsUpdate" json:"teamsUpdate,omitempty"`
	PlayersCreate *Players `protobuf:"bytes,3,opt,name=playersCreate" json:"playersCreate,omitempty"`
	PlayersUpdate *Players `protobuf:"bytes,4,opt,name=playersUpdate" json:"playersUpdate,omitempty"`
	GamesCreate   *Games   `protobuf:"bytes,5,opt,name=gamesCreate" json:"gamesCreate,omitempty"`
	GamesUpdate   *Games   `protobuf:"bytes,6,opt,name=gamesUpdate" json:"gamesUpdate,omitempty"`
}

func (m *BatchUpdates) Reset()                    { *m = BatchUpdates{} }
func (m *BatchUpdates) String() string            { return proto1.CompactTextString(m) }
func (*BatchUpdates) ProtoMessage()               {}
func (*BatchUpdates) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *BatchUpdates) GetTeamsCreate() *Teams {
	if m != nil {
		return m.TeamsCreate
	}
	return nil
}

func (m *BatchUpdates) GetTeamsUpdate() *Teams {
	if m != nil {
		return m.TeamsUpdate
	}
	return nil
}

func (m *BatchUpdates) GetPlayersCreate() *Players {
	if m != nil {
		return m.PlayersCreate
	}
	return nil
}

func (m *BatchUpdates) GetPlayersUpdate() *Players {
	if m != nil {
		return m.PlayersUpdate
	}
	return nil
}

func (m *BatchUpdates) GetGamesCreate() *Games {
	if m != nil {
		return m.GamesCreate
	}
	return nil
}

func (m *BatchUpdates) GetGamesUpdate() *Games {
	if m != nil {
		return m.GamesUpdate
	}
	return nil
}

type Result struct {
	Status Result_Status `protobuf:"varint,1,opt,name=status,enum=proto.Result_Status" json:"status,omitempty"`
	Detail string        `protobuf:"bytes,2,opt,name=detail" json:"detail,omitempty"`
}

func (m *Result) Reset()                    { *m = Result{} }
func (m *Result) String() string            { return proto1.CompactTextString(m) }
func (*Result) ProtoMessage()               {}
func (*Result) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func init() {
	proto1.RegisterType((*Team)(nil), "proto.Team")
	proto1.RegisterType((*Teams)(nil), "proto.Teams")
	proto1.RegisterType((*Player)(nil), "proto.Player")
	proto1.RegisterType((*Players)(nil), "proto.Players")
	proto1.RegisterType((*Game)(nil), "proto.Game")
	proto1.RegisterType((*Games)(nil), "proto.Games")
	proto1.RegisterType((*BatchUpdates)(nil), "proto.BatchUpdates")
	proto1.RegisterType((*Result)(nil), "proto.Result")
	proto1.RegisterEnum("proto.Result_Status", Result_Status_name, Result_Status_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion2

// Client API for ContextUpdate service

type ContextUpdateClient interface {
	CreateTeam(ctx context.Context, in *Team, opts ...grpc.CallOption) (*Result, error)
	UpdateTeam(ctx context.Context, in *Team, opts ...grpc.CallOption) (*Result, error)
	CreatePlayer(ctx context.Context, in *Player, opts ...grpc.CallOption) (*Result, error)
	UpdatePlayer(ctx context.Context, in *Player, opts ...grpc.CallOption) (*Result, error)
	CreateGame(ctx context.Context, in *Game, opts ...grpc.CallOption) (*Result, error)
	UpdateGame(ctx context.Context, in *Game, opts ...grpc.CallOption) (*Result, error)
	BatchUpdate(ctx context.Context, in *BatchUpdates, opts ...grpc.CallOption) (*Result, error)
}

type contextUpdateClient struct {
	cc *grpc.ClientConn
}

func NewContextUpdateClient(cc *grpc.ClientConn) ContextUpdateClient {
	return &contextUpdateClient{cc}
}

func (c *contextUpdateClient) CreateTeam(ctx context.Context, in *Team, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/proto.ContextUpdate/CreateTeam", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contextUpdateClient) UpdateTeam(ctx context.Context, in *Team, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/proto.ContextUpdate/UpdateTeam", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contextUpdateClient) CreatePlayer(ctx context.Context, in *Player, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/proto.ContextUpdate/CreatePlayer", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contextUpdateClient) UpdatePlayer(ctx context.Context, in *Player, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/proto.ContextUpdate/UpdatePlayer", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contextUpdateClient) CreateGame(ctx context.Context, in *Game, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/proto.ContextUpdate/CreateGame", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contextUpdateClient) UpdateGame(ctx context.Context, in *Game, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/proto.ContextUpdate/UpdateGame", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contextUpdateClient) BatchUpdate(ctx context.Context, in *BatchUpdates, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/proto.ContextUpdate/BatchUpdate", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ContextUpdate service

type ContextUpdateServer interface {
	CreateTeam(context.Context, *Team) (*Result, error)
	UpdateTeam(context.Context, *Team) (*Result, error)
	CreatePlayer(context.Context, *Player) (*Result, error)
	UpdatePlayer(context.Context, *Player) (*Result, error)
	CreateGame(context.Context, *Game) (*Result, error)
	UpdateGame(context.Context, *Game) (*Result, error)
	BatchUpdate(context.Context, *BatchUpdates) (*Result, error)
}

func RegisterContextUpdateServer(s *grpc.Server, srv ContextUpdateServer) {
	s.RegisterService(&_ContextUpdate_serviceDesc, srv)
}

func _ContextUpdate_CreateTeam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Team)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextUpdateServer).CreateTeam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ContextUpdate/CreateTeam",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextUpdateServer).CreateTeam(ctx, req.(*Team))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContextUpdate_UpdateTeam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Team)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextUpdateServer).UpdateTeam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ContextUpdate/UpdateTeam",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextUpdateServer).UpdateTeam(ctx, req.(*Team))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContextUpdate_CreatePlayer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Player)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextUpdateServer).CreatePlayer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ContextUpdate/CreatePlayer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextUpdateServer).CreatePlayer(ctx, req.(*Player))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContextUpdate_UpdatePlayer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Player)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextUpdateServer).UpdatePlayer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ContextUpdate/UpdatePlayer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextUpdateServer).UpdatePlayer(ctx, req.(*Player))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContextUpdate_CreateGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Game)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextUpdateServer).CreateGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ContextUpdate/CreateGame",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextUpdateServer).CreateGame(ctx, req.(*Game))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContextUpdate_UpdateGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Game)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextUpdateServer).UpdateGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ContextUpdate/UpdateGame",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextUpdateServer).UpdateGame(ctx, req.(*Game))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContextUpdate_BatchUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchUpdates)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContextUpdateServer).BatchUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.ContextUpdate/BatchUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContextUpdateServer).BatchUpdate(ctx, req.(*BatchUpdates))
	}
	return interceptor(ctx, in, info, handler)
}

var _ContextUpdate_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ContextUpdate",
	HandlerType: (*ContextUpdateServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateTeam",
			Handler:    _ContextUpdate_CreateTeam_Handler,
		},
		{
			MethodName: "UpdateTeam",
			Handler:    _ContextUpdate_UpdateTeam_Handler,
		},
		{
			MethodName: "CreatePlayer",
			Handler:    _ContextUpdate_CreatePlayer_Handler,
		},
		{
			MethodName: "UpdatePlayer",
			Handler:    _ContextUpdate_UpdatePlayer_Handler,
		},
		{
			MethodName: "CreateGame",
			Handler:    _ContextUpdate_CreateGame_Handler,
		},
		{
			MethodName: "UpdateGame",
			Handler:    _ContextUpdate_UpdateGame_Handler,
		},
		{
			MethodName: "BatchUpdate",
			Handler:    _ContextUpdate_BatchUpdate_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 545 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xac, 0x54, 0x5d, 0x8e, 0xd3, 0x30,
	0x10, 0xde, 0x36, 0x7f, 0xed, 0xa4, 0x5d, 0x55, 0xa6, 0x42, 0xd1, 0x0a, 0xc1, 0x92, 0x17, 0x56,
	0x2b, 0x94, 0x87, 0x2c, 0x17, 0x80, 0x0a, 0x10, 0x12, 0x48, 0xc8, 0x65, 0x0f, 0xe0, 0x6d, 0xac,
	0x52, 0x29, 0x6d, 0xaa, 0xc4, 0x95, 0xe0, 0x85, 0x83, 0x71, 0x0e, 0x8e, 0xc1, 0x21, 0xb0, 0x67,
	0xdc, 0xd4, 0x81, 0xdd, 0xaa, 0x0f, 0x3c, 0xc5, 0x33, 0xdf, 0x37, 0xfe, 0x3e, 0xcf, 0xd8, 0x81,
	0xe9, 0xa2, 0xda, 0x28, 0xf9, 0x4d, 0x7d, 0x12, 0x1b, 0xb1, 0x94, 0x75, 0xb6, 0xad, 0x2b, 0x55,
	0xb1, 0x00, 0x3f, 0x69, 0x0d, 0xfe, 0x17, 0x29, 0xd6, 0xec, 0x31, 0x84, 0x4a, 0x7f, 0x57, 0x45,
	0xd2, 0xbb, 0xec, 0x5d, 0x05, 0xdc, 0x46, 0x6c, 0x0a, 0x41, 0xb9, 0x68, 0x74, 0xba, 0x8f, 0x69,
	0x0a, 0x0c, 0xbb, 0x5e, 0x55, 0x4a, 0xa7, 0x3d, 0x62, 0x53, 0xc4, 0x18, 0xf8, 0x1b, 0xb1, 0x96,
	0x89, 0xaf, 0xb3, 0x43, 0x8e, 0x6b, 0x36, 0x01, 0x4f, 0x89, 0x65, 0x12, 0x60, 0xca, 0x2c, 0xd3,
	0x6b, 0x08, 0x8c, 0x66, 0xc3, 0x9e, 0x43, 0x60, 0x64, 0x1a, 0xad, 0xe9, 0x5d, 0xc5, 0x79, 0x4c,
	0xd6, 0x32, 0x03, 0x72, 0x42, 0xd2, 0x1f, 0x10, 0x7e, 0x2e, 0xc5, 0x77, 0x59, 0xb3, 0x0b, 0x18,
	0x6c, 0x71, 0xd5, 0x7a, 0x6c, 0xe3, 0xff, 0xe0, 0xf2, 0x70, 0xfe, 0xc0, 0x3d, 0x7f, 0x9a, 0x43,
	0x44, 0xfa, 0x0d, 0x7b, 0x01, 0x11, 0x09, 0xee, 0xfd, 0x8e, 0xad, 0x5f, 0x22, 0xf0, 0x3d, 0x9a,
	0xfe, 0xee, 0x81, 0xff, 0xde, 0x6e, 0xba, 0xd4, 0xdf, 0x43, 0x53, 0x29, 0x7a, 0xc0, 0xee, 0x53,
	0x00, 0x63, 0xd0, 0x56, 0x78, 0x68, 0xce, 0xc9, 0xb0, 0x4b, 0x88, 0x4d, 0xb4, 0x16, 0x6a, 0xf1,
	0x55, 0x13, 0xc8, 0xbd, 0x9b, 0x62, 0x4f, 0x60, 0x58, 0xcb, 0xa2, 0x73, 0x8e, 0x43, 0xc2, 0xec,
	0x7f, 0x57, 0xee, 0xa4, 0x85, 0x43, 0x84, 0x9d, 0x8c, 0xa9, 0x36, 0x4a, 0x8d, 0x12, 0xb5, 0x4a,
	0x22, 0x0d, 0x7b, 0xfc, 0x90, 0x60, 0x09, 0x44, 0x26, 0x90, 0x9b, 0x22, 0x19, 0x20, 0xb6, 0x0f,
	0xcd, 0x38, 0xcd, 0x69, 0x71, 0x9c, 0xc8, 0xff, 0x6b, 0x9c, 0x06, 0xe4, 0x84, 0xa4, 0x3f, 0xfb,
	0x30, 0x7a, 0x63, 0xdc, 0xde, 0x6e, 0x0b, 0xa1, 0x74, 0x4d, 0x06, 0x31, 0x0e, 0x7a, 0x56, 0x4b,
	0x1d, 0x63, 0x9f, 0xe2, 0x7c, 0xe4, 0x5c, 0x84, 0x86, 0xbb, 0x84, 0x96, 0x4f, 0xf5, 0xd8, 0xc0,
	0xfb, 0xf9, 0x44, 0x60, 0xaf, 0x60, 0x6c, 0xc7, 0x62, 0x15, 0x3c, 0xac, 0x38, 0xef, 0x8c, 0xae,
	0xe1, 0x5d, 0x92, 0x53, 0x65, 0x75, 0xfc, 0xa3, 0x55, 0x56, 0x4b, 0x7b, 0xc3, 0x53, 0x5a, 0xa5,
	0xa0, 0xe3, 0x0d, 0x5b, 0xc4, 0x5d, 0x42, 0xcb, 0xb7, 0x1a, 0xe1, 0x83, 0x7c, 0x22, 0xa4, 0x15,
	0x84, 0x5c, 0x36, 0xbb, 0x52, 0xb1, 0x97, 0x10, 0xea, 0xa9, 0xa8, 0x5d, 0x83, 0x0d, 0x3b, 0xcf,
	0xa7, 0xb6, 0x88, 0xe0, 0x6c, 0x8e, 0x18, 0xb7, 0x1c, 0x73, 0x0d, 0x0b, 0xa9, 0xc4, 0xaa, 0xc4,
	0x76, 0x0d, 0xb9, 0x8d, 0xd2, 0x67, 0x10, 0x12, 0x93, 0xc5, 0x10, 0xcd, 0x6f, 0x67, 0xb3, 0xb7,
	0xf3, 0xf9, 0xe4, 0x8c, 0x0d, 0xc0, 0x7f, 0xf7, 0xfa, 0xc3, 0xc7, 0x49, 0x2f, 0xff, 0xd5, 0x87,
	0xf1, 0x8c, 0x7e, 0x1e, 0xf6, 0x88, 0xd7, 0x00, 0x64, 0x1e, 0x7f, 0x1a, 0xee, 0x83, 0xbd, 0x18,
	0x77, 0x3c, 0xa4, 0x67, 0x86, 0x4b, 0x55, 0x27, 0x70, 0x33, 0x18, 0xd1, 0xbe, 0xf6, 0xb1, 0x77,
	0x9f, 0xd6, 0xbd, 0x7c, 0xda, 0xfb, 0x44, 0x7e, 0xeb, 0x1b, 0xdf, 0xa5, 0x7b, 0x33, 0x8f, 0xf8,
	0x3e, 0x81, 0x7b, 0x03, 0xb1, 0x73, 0x9d, 0xd9, 0x23, 0x8b, 0xbb, 0x57, 0xfc, 0x9f, 0xa2, 0xbb,
	0x10, 0xe3, 0x9b, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x26, 0xa8, 0x76, 0xfc, 0x99, 0x05, 0x00,
	0x00,
}
