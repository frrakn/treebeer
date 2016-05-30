// Code generated by protoc-gen-go.
// source: context.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	context.proto

It has these top-level messages:
	Empty
	Team
	Player
	Game
	Stat
	SeasonUpdates
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

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto1.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Team struct {
	Lcsid  int32  `protobuf:"varint,2,opt,name=lcsid" json:"lcsid,omitempty"`
	Riotid int32  `protobuf:"varint,3,opt,name=riotid" json:"riotid,omitempty"`
	Name   string `protobuf:"bytes,4,opt,name=name" json:"name,omitempty"`
	Tag    string `protobuf:"bytes,5,opt,name=tag" json:"tag,omitempty"`
}

func (m *Team) Reset()                    { *m = Team{} }
func (m *Team) String() string            { return proto1.CompactTextString(m) }
func (*Team) ProtoMessage()               {}
func (*Team) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type Player struct {
	Lcsid    int32    `protobuf:"varint,2,opt,name=lcsid" json:"lcsid,omitempty"`
	Riotid   int32    `protobuf:"varint,3,opt,name=riotid" json:"riotid,omitempty"`
	Name     string   `protobuf:"bytes,4,opt,name=name" json:"name,omitempty"`
	Teamid   int32    `protobuf:"varint,5,opt,name=teamid" json:"teamid,omitempty"`
	Position string   `protobuf:"bytes,6,opt,name=position" json:"position,omitempty"`
	Addlpos  []string `protobuf:"bytes,7,rep,name=addlpos" json:"addlpos,omitempty"`
}

func (m *Player) Reset()                    { *m = Player{} }
func (m *Player) String() string            { return proto1.CompactTextString(m) }
func (*Player) ProtoMessage()               {}
func (*Player) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type Game struct {
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
func (*Game) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type Stat struct {
	Riotname string `protobuf:"bytes,2,opt,name=riotname" json:"riotname,omitempty"`
}

func (m *Stat) Reset()                    { *m = Stat{} }
func (m *Stat) String() string            { return proto1.CompactTextString(m) }
func (*Stat) ProtoMessage()               {}
func (*Stat) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type SeasonUpdates struct {
	Teams   []*Team   `protobuf:"bytes,1,rep,name=teams" json:"teams,omitempty"`
	Players []*Player `protobuf:"bytes,2,rep,name=players" json:"players,omitempty"`
}

func (m *SeasonUpdates) Reset()                    { *m = SeasonUpdates{} }
func (m *SeasonUpdates) String() string            { return proto1.CompactTextString(m) }
func (*SeasonUpdates) ProtoMessage()               {}
func (*SeasonUpdates) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *SeasonUpdates) GetTeams() []*Team {
	if m != nil {
		return m.Teams
	}
	return nil
}

func (m *SeasonUpdates) GetPlayers() []*Player {
	if m != nil {
		return m.Players
	}
	return nil
}

func init() {
	proto1.RegisterType((*Empty)(nil), "proto.Empty")
	proto1.RegisterType((*Team)(nil), "proto.Team")
	proto1.RegisterType((*Player)(nil), "proto.Player")
	proto1.RegisterType((*Game)(nil), "proto.Game")
	proto1.RegisterType((*Stat)(nil), "proto.Stat")
	proto1.RegisterType((*SeasonUpdates)(nil), "proto.SeasonUpdates")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion2

// Client API for SeasonUpdate service

type SeasonUpdateClient interface {
	SeasonUpdate(ctx context.Context, in *SeasonUpdates, opts ...grpc.CallOption) (*Empty, error)
}

type seasonUpdateClient struct {
	cc *grpc.ClientConn
}

func NewSeasonUpdateClient(cc *grpc.ClientConn) SeasonUpdateClient {
	return &seasonUpdateClient{cc}
}

func (c *seasonUpdateClient) SeasonUpdate(ctx context.Context, in *SeasonUpdates, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/proto.SeasonUpdate/SeasonUpdate", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for SeasonUpdate service

type SeasonUpdateServer interface {
	SeasonUpdate(context.Context, *SeasonUpdates) (*Empty, error)
}

func RegisterSeasonUpdateServer(s *grpc.Server, srv SeasonUpdateServer) {
	s.RegisterService(&_SeasonUpdate_serviceDesc, srv)
}

func _SeasonUpdate_SeasonUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SeasonUpdates)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeasonUpdateServer).SeasonUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.SeasonUpdate/SeasonUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeasonUpdateServer).SeasonUpdate(ctx, req.(*SeasonUpdates))
	}
	return interceptor(ctx, in, info, handler)
}

var _SeasonUpdate_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.SeasonUpdate",
	HandlerType: (*SeasonUpdateServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SeasonUpdate",
			Handler:    _SeasonUpdate_SeasonUpdate_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

// Client API for LiveStatUpdate service

type LiveStatUpdateClient interface {
	GetGame(ctx context.Context, in *Game, opts ...grpc.CallOption) (*Empty, error)
	GetStat(ctx context.Context, in *Stat, opts ...grpc.CallOption) (*Empty, error)
}

type liveStatUpdateClient struct {
	cc *grpc.ClientConn
}

func NewLiveStatUpdateClient(cc *grpc.ClientConn) LiveStatUpdateClient {
	return &liveStatUpdateClient{cc}
}

func (c *liveStatUpdateClient) GetGame(ctx context.Context, in *Game, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/proto.LiveStatUpdate/GetGame", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *liveStatUpdateClient) GetStat(ctx context.Context, in *Stat, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/proto.LiveStatUpdate/GetStat", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for LiveStatUpdate service

type LiveStatUpdateServer interface {
	GetGame(context.Context, *Game) (*Empty, error)
	GetStat(context.Context, *Stat) (*Empty, error)
}

func RegisterLiveStatUpdateServer(s *grpc.Server, srv LiveStatUpdateServer) {
	s.RegisterService(&_LiveStatUpdate_serviceDesc, srv)
}

func _LiveStatUpdate_GetGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Game)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LiveStatUpdateServer).GetGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LiveStatUpdate/GetGame",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LiveStatUpdateServer).GetGame(ctx, req.(*Game))
	}
	return interceptor(ctx, in, info, handler)
}

func _LiveStatUpdate_GetStat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Stat)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LiveStatUpdateServer).GetStat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LiveStatUpdate/GetStat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LiveStatUpdateServer).GetStat(ctx, req.(*Stat))
	}
	return interceptor(ctx, in, info, handler)
}

var _LiveStatUpdate_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.LiveStatUpdate",
	HandlerType: (*LiveStatUpdateServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetGame",
			Handler:    _LiveStatUpdate_GetGame_Handler,
		},
		{
			MethodName: "GetStat",
			Handler:    _LiveStatUpdate_GetStat_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 393 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xac, 0x92, 0xcf, 0x4e, 0xe3, 0x30,
	0x10, 0xc6, 0x37, 0xcd, 0xbf, 0x76, 0xda, 0xae, 0x56, 0x56, 0xb5, 0xb2, 0xaa, 0xd5, 0xaa, 0x9b,
	0xc3, 0xd2, 0x53, 0x0f, 0x85, 0x47, 0x00, 0xf5, 0xc2, 0x01, 0xa5, 0x70, 0x81, 0x0b, 0x6e, 0x63,
	0x95, 0x48, 0xf9, 0xa7, 0x64, 0x40, 0xf4, 0x4d, 0x78, 0x30, 0x1e, 0x08, 0x8f, 0xed, 0xd0, 0xb4,
	0xe2, 0xc8, 0x29, 0xf9, 0x7d, 0xf3, 0xcd, 0xe4, 0xcb, 0xd8, 0x30, 0xde, 0x96, 0x05, 0xca, 0x57,
	0x5c, 0x54, 0x75, 0x89, 0x25, 0xf3, 0xf5, 0x23, 0x0a, 0xc1, 0xbf, 0xca, 0x2b, 0xdc, 0x47, 0xf7,
	0xe0, 0xdd, 0x4a, 0x91, 0xb3, 0x09, 0xf8, 0xd9, 0xb6, 0x49, 0x13, 0xde, 0x9b, 0x39, 0x73, 0x3f,
	0x36, 0xc0, 0x7e, 0x43, 0x50, 0xa7, 0x25, 0x2a, 0xd9, 0xd5, 0xb2, 0x25, 0xc6, 0xc0, 0x2b, 0x44,
	0x2e, 0xb9, 0xa7, 0xd4, 0x41, 0xac, 0xdf, 0xd9, 0x2f, 0x70, 0x51, 0xec, 0xb8, 0xaf, 0x25, 0x7a,
	0x8d, 0xde, 0x1c, 0x08, 0x6e, 0x32, 0xb1, 0x97, 0xf5, 0x37, 0x8c, 0x57, 0x5e, 0x54, 0x41, 0x95,
	0xd7, 0x37, 0x5e, 0x43, 0x6c, 0x0a, 0xfd, 0xaa, 0x6c, 0x52, 0x4c, 0xcb, 0x82, 0x07, 0xda, 0xff,
	0xc9, 0x8c, 0x43, 0x28, 0x92, 0x24, 0x53, 0xcc, 0xc3, 0x99, 0xab, 0x4a, 0x2d, 0x46, 0xef, 0x0e,
	0x78, 0x2b, 0x1a, 0xfb, 0x75, 0xb0, 0xbf, 0x00, 0x14, 0x65, 0xa7, 0x1c, 0x36, 0xdc, 0x20, 0xee,
	0x28, 0x6c, 0x06, 0x43, 0xa2, 0x5c, 0xe0, 0xf6, 0x49, 0x19, 0x4c, 0xce, 0xae, 0xc4, 0xfe, 0xc0,
	0xa0, 0x96, 0xc9, 0x51, 0xe2, 0x83, 0x40, 0xf3, 0x37, 0xd9, 0xb3, 0xb4, 0xe5, 0x40, 0x97, 0x3b,
	0x0a, 0x75, 0xd3, 0x97, 0x1a, 0x14, 0x35, 0xaa, 0xe8, 0xce, 0xdc, 0x8d, 0x0f, 0x02, 0xfd, 0x16,
	0x81, 0x2c, 0x12, 0xde, 0xd7, 0xb5, 0x16, 0xa3, 0x08, 0xbc, 0x35, 0x0a, 0xa4, 0xa5, 0x50, 0x18,
	0xbd, 0xc4, 0x9e, 0x59, 0x4a, 0xcb, 0xd1, 0x03, 0x8c, 0xd7, 0x52, 0x34, 0x65, 0x71, 0x57, 0x25,
	0x02, 0x65, 0xc3, 0xfe, 0x81, 0x4f, 0x9f, 0x6d, 0xb8, 0xa3, 0x76, 0x34, 0x5c, 0x0e, 0xcd, 0x4d,
	0x59, 0xd0, 0xb5, 0x88, 0x4d, 0x85, 0x9d, 0x41, 0x58, 0xe9, 0x83, 0x6c, 0xd4, 0x38, 0x32, 0x8d,
	0xad, 0xc9, 0x1c, 0x6f, 0xdc, 0x56, 0x97, 0x97, 0x30, 0xea, 0x0e, 0x67, 0x17, 0x27, 0x3c, 0xb1,
	0x7d, 0x47, 0x09, 0xa6, 0x23, 0xab, 0x9a, 0x2b, 0xf9, 0x63, 0xf9, 0x08, 0x3f, 0xaf, 0xd3, 0x17,
	0x49, 0xbf, 0x62, 0xfb, 0xfe, 0x43, 0xb8, 0x92, 0xa8, 0x4f, 0xac, 0xcd, 0x47, 0x70, 0xda, 0x69,
	0x7d, 0x7a, 0x07, 0xad, 0x8f, 0xe0, 0xd4, 0xb7, 0x09, 0x34, 0x9e, 0x7f, 0x04, 0x00, 0x00, 0xff,
	0xff, 0x6e, 0x3c, 0x19, 0xcb, 0x1e, 0x03, 0x00, 0x00,
}
