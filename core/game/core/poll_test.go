package core

// import (
// 	"fmt"
// 	"testing"
// 	"uwwolf/game/enum"
// 	"uwwolf/game/types"

// 	"github.com/stretchr/testify/suite"
// )

// type PollSuite struct {
// 	suite.Suite
// }

// func TestPollSuite(t *testing.T) {
// 	suite.Run(t, new(PollSuite))
// }

// func (ps *PollSuite) TestNewPoll() {
// 	tests := []struct {
// 		name        string
// 		capacity    uint8
// 		expectedErr string
// 	}{
// 		{
// 			name:        "Failure (Too small capacity)",
// 			capacity:    0,
// 			expectedErr: "The capacity is too small ¬_¬",
// 		},
// 		{
// 			name:     "Ok",
// 			capacity: 5,
// 		},
// 	}

// 	for _, test := range tests {
// 		ps.Run(test.name, func() {
// 			myPoll, err := NewPoll(test.capacity)

// 			if test.expectedErr != "" {
// 				ps.NotNil(err)
// 				ps.Equal(test.expectedErr, err.Error())
// 			} else {
// 				ps.NotNil(myPoll)
// 				ps.NotNil(myPoll.(*poll).Weights)
// 				ps.NotNil(myPoll.(*poll).Records)
// 				ps.Equal(test.capacity, myPoll.(*poll).Capacity)
// 			}
// 		})
// 	}
// }

// func (ps *PollSuite) TestIsOpen() {
// 	tests := []struct {
// 		name           string
// 		expectedStatus bool
// 		expectedRound  enum.Round
// 		setup          func(*poll)
// 	}{
// 		{
// 			name:           "Not open (Round is zero)",
// 			expectedStatus: false,
// 			expectedRound:  0,
// 			setup: func(myPoll *poll) {
// 				myPoll.Round = 0
// 			},
// 		},
// 		{
// 			name:           "Not open (Poll record is nil)",
// 			expectedStatus: false,
// 			expectedRound:  1,
// 			setup: func(myPoll *poll) {
// 				myPoll.Round = 1
// 			},
// 		},
// 		{
// 			name:           "Not open (Poll was closed)",
// 			expectedStatus: false,
// 			expectedRound:  1,
// 			setup: func(myPoll *poll) {
// 				myPoll.Round = 1
// 				myPoll.Records[myPoll.Round] = &types.PollRecord{
// 					IsClosed: true,
// 				}
// 			},
// 		},
// 		{
// 			name:           "Is open",
// 			expectedStatus: true,
// 			expectedRound:  1,
// 			setup: func(myPoll *poll) {
// 				myPoll.Round = 1
// 				myPoll.Records[myPoll.Round] = &types.PollRecord{
// 					IsClosed: false,
// 				}
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ps.Run(test.name, func() {
// 			myPoll, _ := NewPoll(9)
// 			test.setup(myPoll.(*poll))
// 			isOpen, round := myPoll.IsOpen()

// 			ps.Equal(test.expectedStatus, isOpen)
// 			ps.Equal(test.expectedRound, round)
// 		})
// 	}
// }

// func (ps *PollSuite) TestCanVote() {
// 	electorID := enum.PlayerID("1")
// 	tests := []struct {
// 		name           string
// 		electorID      enum.PlayerID
// 		expectedStatus bool
// 		expectedErr    string
// 		setup          func(*poll)
// 	}{
// 		{
// 			name:           "Cannot vote (Poll was closed)",
// 			electorID:      electorID,
// 			expectedStatus: false,
// 			expectedErr:    fmt.Sprintf("Poll (%v) is closed ᕙ(⇀‸↼‶)ᕗ", 0),
// 			setup: func(myPoll *poll) {
// 				myPoll.Round = 0
// 			},
// 		},
// 		{
// 			name:           "Cannot vote (Not an elector)",
// 			electorID:      electorID,
// 			expectedStatus: false,
// 			expectedErr:    "You're not allowed to vote ノ(ジ)ー'",
// 			setup: func(myPoll *poll) {
// 				// Open poll
// 				myPoll.Round = 1
// 				myPoll.Records[myPoll.Round] = &types.PollRecord{
// 					IsClosed: false,
// 				}
// 			},
// 		},
// 		{
// 			name:           "Cannot vote (Already voted)",
// 			electorID:      electorID,
// 			expectedStatus: false,
// 			expectedErr:    "Wait for the next round ಠ_ಠ",
// 			setup: func(myPoll *poll) {
// 				// Open poll
// 				myPoll.Round = 1
// 				myPoll.Records[myPoll.Round] = &types.PollRecord{
// 					IsClosed: false,
// 				}

// 				myPoll.RemainingElectorIDs = append(
// 					myPoll.RemainingElectorIDs,
// 					electorID,
// 				)
// 				myPoll.VotedElectorIDs = append(
// 					myPoll.VotedElectorIDs,
// 					electorID,
// 				)
// 			},
// 		},
// 		{
// 			name:           "Ok",
// 			electorID:      electorID,
// 			expectedStatus: true,
// 			setup: func(myPoll *poll) {
// 				// Open poll
// 				myPoll.Round = 1
// 				myPoll.Records[myPoll.Round] = &types.PollRecord{
// 					IsClosed: false,
// 				}

// 				myPoll.RemainingElectorIDs = append(
// 					myPoll.RemainingElectorIDs,
// 					electorID,
// 				)
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ps.Run(test.name, func() {
// 			myPoll, _ := NewPoll(9)
// 			test.setup(myPoll.(*poll))
// 			res, err := myPoll.CanVote(test.electorID)

// 			ps.Equal(test.expectedStatus, res)

// 			if res == false {
// 				ps.Equal(test.expectedErr, err.Error())
// 			}
// 		})
// 	}
// }

// func (ps *PollSuite) TestRecord() {
// 	tests := []struct {
// 		name           string
// 		round          enum.Round
// 		expectedRecord *types.PollRecord
// 		setup          func(*poll)
// 	}{
// 		{
// 			name:           "Nil (Round is zero)",
// 			round:          0,
// 			expectedRecord: nil,
// 			setup:          func(myPoll *poll) {},
// 		},
// 		{
// 			name:           "Nil (Non-existent round)",
// 			round:          99,
// 			expectedRecord: nil,
// 			setup: func(myPoll *poll) {
// 				myPoll.Round = 1
// 				myPoll.Records[myPoll.Round] = &types.PollRecord{}
// 			},
// 		},
// 		{
// 			name:  "Ok",
// 			round: 1,
// 			expectedRecord: &types.PollRecord{
// 				WinnerID: "98",
// 			},
// 			setup: func(myPoll *poll) {
// 				myPoll.Round = 1
// 				myPoll.Records[myPoll.Round] = &types.PollRecord{
// 					WinnerID: "98",
// 				}
// 			},
// 		},
// 		{
// 			name:  "Ok (Get latest record)",
// 			round: enum.LastRound,
// 			expectedRecord: &types.PollRecord{
// 				WinnerID: "99",
// 			},
// 			setup: func(myPoll *poll) {
// 				myPoll.Records[1] = &types.PollRecord{
// 					WinnerID: "98",
// 				}
// 				myPoll.Round = 2
// 				myPoll.Records[myPoll.Round] = &types.PollRecord{
// 					WinnerID: "99",
// 				}
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ps.Run(test.name, func() {
// 			myPoll, _ := NewPoll(9)
// 			test.setup(myPoll.(*poll))
// 			record := myPoll.Record(test.round)

// 			ps.Equal(test.expectedRecord, record)
// 		})
// 	}
// }

// func (ps *PollSuite) TestOpen() {
// 	capacity := uint8(5)
// 	tests := []struct {
// 		name           string
// 		expectedStatus bool
// 		expectedRound  enum.Round
// 		setup          func(*poll)
// 	}{
// 		{
// 			name:           "Failure (Already open)",
// 			expectedStatus: false,
// 			expectedRound:  1,
// 			setup: func(myPoll *poll) {
// 				myPoll.Round = 1
// 				myPoll.Records[myPoll.Round] = &types.PollRecord{
// 					IsClosed: false,
// 				}
// 			},
// 		},
// 		{
// 			name:           "Failure (Not enough electors)",
// 			expectedStatus: false,
// 			expectedRound:  1,
// 			setup: func(myPoll *poll) {
// 				myPoll.Round = 1
// 				myPoll.Records[myPoll.Round] = &types.PollRecord{
// 					IsClosed: true,
// 				}
// 				myPoll.RemainingElectorIDs = []enum.PlayerID{"1", "2", "3"}
// 			},
// 		},
// 		{
// 			name:           "Ok",
// 			expectedStatus: true,
// 			expectedRound:  2,
// 			setup: func(myPoll *poll) {
// 				myPoll.Round = 1
// 				myPoll.Records[myPoll.Round] = &types.PollRecord{
// 					IsClosed: true,
// 				}
// 				myPoll.VotedElectorIDs = []enum.PlayerID{"1", "2", "3", "4", "5"}
// 				myPoll.RemainingElectorIDs = []enum.PlayerID{"1", "2", "3", "4", "5"}
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ps.Run(test.name, func() {
// 			myPoll, _ := NewPoll(capacity)
// 			test.setup(myPoll.(*poll))
// 			ok, round := myPoll.Open()

// 			ps.Equal(test.expectedStatus, ok)
// 			ps.Equal(test.expectedRound, round)

// 			if test.expectedStatus == true {
// 				ps.Empty(myPoll.(*poll).VotedElectorIDs)
// 				ps.False(myPoll.Record(enum.LastRound).IsClosed)
// 			}
// 		})
// 	}
// }

// func (ps *PollSuite) TestClose() {
// 	tests := []struct {
// 		name           string
// 		expectedStatus bool
// 		expectedRecord *types.PollRecord
// 		setup          func(*poll)
// 	}{
// 		{
// 			name:           "Failure (Poll was closed)",
// 			expectedStatus: false,
// 			expectedRecord: nil,
// 			setup: func(myPoll *poll) {
// 				myPoll.Round = 1
// 				myPoll.Records[myPoll.Round] = &types.PollRecord{
// 					IsClosed: true,
// 				}
// 			},
// 		},
// 		{
// 			name:           "Ok (No one have half-votes weight)",
// 			expectedStatus: true,
// 			expectedRecord: &types.PollRecord{
// 				IsClosed: true,
// 				WinnerID: "",
// 				VoteRecords: map[enum.PlayerID]*types.VoteRecord{
// 					"1": {
// 						ElectorIDs: []enum.PlayerID{"2", "3"},
// 						Weights:    2,
// 						Votes:      2,
// 					},
// 					"2": {
// 						ElectorIDs: []enum.PlayerID{"1", "4"},
// 						Weights:    2,
// 						Votes:      2,
// 					},
// 					"": {
// 						ElectorIDs: []enum.PlayerID{"5"},
// 						Votes:      1,
// 						Weights:    1,
// 					},
// 				},
// 			},
// 			setup: func(myPoll *poll) {
// 				myPoll.Round = 1
// 				myPoll.RemainingElectorIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 					"3",
// 					"4",
// 					"5",
// 				}
// 				myPoll.VotedElectorIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 					"3",
// 					"4",
// 				}
// 				myPoll.Records[myPoll.Round] = &types.PollRecord{
// 					IsClosed: false,
// 					VoteRecords: map[enum.PlayerID]*types.VoteRecord{
// 						"1": {
// 							ElectorIDs: []enum.PlayerID{"2", "3"},
// 							Weights:    2,
// 							Votes:      2,
// 						},
// 						"2": {
// 							ElectorIDs: []enum.PlayerID{"1", "4"},
// 							Weights:    2,
// 							Votes:      2,
// 						},
// 					},
// 				}
// 				myPoll.Records[myPoll.Round].VoteRecords[""] = &types.VoteRecord{
// 					ElectorIDs: []enum.PlayerID{},
// 				}
// 			},
// 		},
// 		{
// 			name:           "Ok (A draw)",
// 			expectedStatus: true,
// 			expectedRecord: &types.PollRecord{
// 				IsClosed: true,
// 				WinnerID: "",
// 				VoteRecords: map[enum.PlayerID]*types.VoteRecord{
// 					"1": {
// 						ElectorIDs: []enum.PlayerID{"2", "3", "4"},
// 						Weights:    3,
// 						Votes:      3,
// 					},
// 					"2": {
// 						ElectorIDs: []enum.PlayerID{"1", "5", "6"},
// 						Weights:    3,
// 						Votes:      3,
// 					},
// 					"": {
// 						ElectorIDs: []enum.PlayerID{},
// 						Votes:      0,
// 						Weights:    0,
// 					},
// 				},
// 			},
// 			setup: func(myPoll *poll) {
// 				myPoll.Round = 1
// 				myPoll.RemainingElectorIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 					"3",
// 					"4",
// 					"5",
// 					"6",
// 				}
// 				myPoll.VotedElectorIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 					"3",
// 					"4",
// 					"5",
// 					"6",
// 				}
// 				myPoll.Records[myPoll.Round] = &types.PollRecord{
// 					IsClosed: false,
// 					VoteRecords: map[enum.PlayerID]*types.VoteRecord{
// 						"1": {
// 							ElectorIDs: []enum.PlayerID{"2", "3", "4"},
// 							Weights:    3,
// 							Votes:      3,
// 						},
// 						"2": {
// 							ElectorIDs: []enum.PlayerID{"1", "5", "6"},
// 							Weights:    3,
// 							Votes:      3,
// 						},
// 					},
// 				}
// 				myPoll.Records[myPoll.Round].VoteRecords[""] = &types.VoteRecord{
// 					ElectorIDs: []enum.PlayerID{},
// 				}
// 			},
// 		},
// 		{
// 			name:           "Ok (One cadidate have half-votes weight)",
// 			expectedStatus: true,
// 			expectedRecord: &types.PollRecord{
// 				IsClosed: true,
// 				WinnerID: "1",
// 				VoteRecords: map[enum.PlayerID]*types.VoteRecord{
// 					"1": {
// 						ElectorIDs: []enum.PlayerID{"2", "3", "4"},
// 						Weights:    3,
// 						Votes:      3,
// 					},
// 					"2": {
// 						ElectorIDs: []enum.PlayerID{"1"},
// 						Weights:    1,
// 						Votes:      1,
// 					},
// 					"": {
// 						ElectorIDs: []enum.PlayerID{"5"},
// 						Votes:      1,
// 						Weights:    1,
// 					},
// 				},
// 			},
// 			setup: func(myPoll *poll) {
// 				myPoll.Round = 1
// 				myPoll.RemainingElectorIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 					"3",
// 					"4",
// 					"5",
// 				}
// 				myPoll.VotedElectorIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 					"3",
// 					"4",
// 				}
// 				myPoll.Records[myPoll.Round] = &types.PollRecord{
// 					IsClosed: false,
// 					VoteRecords: map[enum.PlayerID]*types.VoteRecord{
// 						"1": {
// 							ElectorIDs: []enum.PlayerID{"2", "3", "4"},
// 							Weights:    3,
// 							Votes:      3,
// 						},
// 						"2": {
// 							ElectorIDs: []enum.PlayerID{"1"},
// 							Weights:    1,
// 							Votes:      1,
// 						},
// 					},
// 				}
// 				myPoll.Records[myPoll.Round].VoteRecords[""] = &types.VoteRecord{
// 					ElectorIDs: []enum.PlayerID{},
// 				}
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ps.Run(test.name, func() {
// 			myPoll, _ := NewPoll(9)
// 			test.setup(myPoll.(*poll))
// 			ok, record := myPoll.Close()

// 			ps.Equal(test.expectedStatus, ok)
// 			ps.Equal(test.expectedRecord, record)

// 			if test.expectedStatus == true {
// 				ps.Equal(
// 					len(myPoll.(*poll).RemainingElectorIDs),
// 					len(myPoll.(*poll).VotedElectorIDs),
// 				)

// 				isOpen, _ := myPoll.IsOpen()
// 				ps.False(isOpen)
// 			}
// 		})
// 	}
// }

// func (ps *PollSuite) TestAddCandidates() {
// 	tests := []struct {
// 		name                     string
// 		candidateIDs             []enum.PlayerID
// 		newCandidateIDs          []enum.PlayerID
// 		newRemainingCandidateIDs []enum.PlayerID
// 		setup                    func(*poll)
// 	}{
// 		{
// 			name: "Failure (Already existed in remaining)",
// 			candidateIDs: []enum.PlayerID{
// 				"1",
// 				"2",
// 			},
// 			newCandidateIDs: []enum.PlayerID{
// 				"1",
// 				"2",
// 				"3",
// 			},
// 			newRemainingCandidateIDs: []enum.PlayerID{
// 				"1",
// 				"2",
// 				"3",
// 			},
// 			setup: func(myPoll *poll) {
// 				myPoll.CandidateIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 					"3",
// 				}
// 				myPoll.RemainingCandidateIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 					"3",
// 				}
// 			},
// 		},
// 		{
// 			name: "Ok (Doesn't exist in remaining but in all)",
// 			candidateIDs: []enum.PlayerID{
// 				"1",
// 				"2",
// 			},
// 			newCandidateIDs: []enum.PlayerID{
// 				"1",
// 				"2",
// 				"3",
// 			},
// 			newRemainingCandidateIDs: []enum.PlayerID{
// 				"2",
// 				"3",
// 				"1",
// 			},
// 			setup: func(myPoll *poll) {
// 				myPoll.CandidateIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 					"3",
// 				}
// 				myPoll.RemainingCandidateIDs = []enum.PlayerID{
// 					"2",
// 					"3",
// 				}
// 			},
// 		},
// 		{
// 			name: "Ok (Doesn't exist in remaining and all)",
// 			candidateIDs: []enum.PlayerID{
// 				"1",
// 				"2",
// 			},
// 			newCandidateIDs: []enum.PlayerID{
// 				"2",
// 				"3",
// 				"1",
// 			},
// 			newRemainingCandidateIDs: []enum.PlayerID{
// 				"2",
// 				"3",
// 				"1",
// 			},
// 			setup: func(myPoll *poll) {
// 				myPoll.CandidateIDs = []enum.PlayerID{
// 					"2",
// 					"3",
// 				}
// 				myPoll.RemainingCandidateIDs = []enum.PlayerID{
// 					"2",
// 					"3",
// 				}
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ps.Run(test.name, func() {
// 			myPoll, _ := NewPoll(9)
// 			test.setup(myPoll.(*poll))
// 			myPoll.AddCandidates(test.candidateIDs...)

// 			ps.Equal(test.newCandidateIDs, myPoll.(*poll).CandidateIDs)
// 			ps.Equal(test.newRemainingCandidateIDs, myPoll.(*poll).RemainingCandidateIDs)
// 		})
// 	}
// }

// func (ps *PollSuite) TestRemoveCandidate() {
// 	tests := []struct {
// 		name                     string
// 		candidateID              enum.PlayerID
// 		expectedStatus           bool
// 		newCandidateIDs          []enum.PlayerID
// 		newRemainingCandidateIDs []enum.PlayerID
// 		setup                    func(*poll)
// 	}{
// 		{
// 			name:           "Failure (Non-existent candidate)",
// 			candidateID:    "99",
// 			expectedStatus: false,
// 			newCandidateIDs: []enum.PlayerID{
// 				"1",
// 				"2",
// 				"3",
// 			},
// 			newRemainingCandidateIDs: []enum.PlayerID{
// 				"1",
// 				"2",
// 				"3",
// 			},
// 			setup: func(myPoll *poll) {
// 				myPoll.CandidateIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 					"3",
// 				}
// 				myPoll.RemainingCandidateIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 					"3",
// 				}
// 			},
// 		},
// 		{
// 			name:           "Ok",
// 			candidateID:    "1",
// 			expectedStatus: true,
// 			newCandidateIDs: []enum.PlayerID{
// 				"1",
// 				"2",
// 				"3",
// 			},
// 			newRemainingCandidateIDs: []enum.PlayerID{
// 				"2",
// 				"3",
// 			},
// 			setup: func(myPoll *poll) {
// 				myPoll.CandidateIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 					"3",
// 				}
// 				myPoll.RemainingCandidateIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 					"3",
// 				}
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ps.Run(test.name, func() {
// 			myPoll, _ := NewPoll(9)
// 			test.setup(myPoll.(*poll))
// 			ok := myPoll.RemoveCandidate(test.candidateID)

// 			ps.Equal(test.expectedStatus, ok)
// 			ps.Equal(test.newCandidateIDs, myPoll.(*poll).CandidateIDs)
// 			ps.Equal(test.newRemainingCandidateIDs, myPoll.(*poll).RemainingCandidateIDs)
// 		})
// 	}
// }

// func (ps *PollSuite) TestAddElectors() {
// 	capacity := uint8(5)
// 	tests := []struct {
// 		name                   string
// 		electorIDs             []enum.PlayerID
// 		expectedStatus         bool
// 		newElectorIDs          []enum.PlayerID
// 		newRemainingElectorIDs []enum.PlayerID
// 		setup                  func(*poll)
// 	}{
// 		{
// 			name: "Failure (Overload)",
// 			electorIDs: []enum.PlayerID{
// 				"4",
// 				"5",
// 				"7",
// 				"8",
// 			},
// 			expectedStatus: false,
// 			newElectorIDs: []enum.PlayerID{
// 				"1",
// 				"2",
// 				"3",
// 			},
// 			newRemainingElectorIDs: []enum.PlayerID{
// 				"1",
// 				"2",
// 			},
// 			setup: func(myPoll *poll) {
// 				myPoll.ElectorIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 					"3",
// 				}
// 				myPoll.RemainingElectorIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 				}
// 			},
// 		},
// 		{
// 			name: "Failure (Already existed in remaining)",
// 			electorIDs: []enum.PlayerID{
// 				"1",
// 				"2",
// 			},
// 			expectedStatus: true,
// 			newElectorIDs: []enum.PlayerID{
// 				"1",
// 				"2",
// 				"3",
// 			},
// 			newRemainingElectorIDs: []enum.PlayerID{
// 				"1",
// 				"2",
// 				"3",
// 			},
// 			setup: func(myPoll *poll) {
// 				myPoll.ElectorIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 					"3",
// 				}
// 				myPoll.RemainingElectorIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 					"3",
// 				}
// 			},
// 		},
// 		{
// 			name: "Ok (Doesn't exist in remaining but in all)",
// 			electorIDs: []enum.PlayerID{
// 				"1",
// 				"2",
// 			},
// 			expectedStatus: true,
// 			newElectorIDs: []enum.PlayerID{
// 				"1",
// 				"2",
// 				"3",
// 			},
// 			newRemainingElectorIDs: []enum.PlayerID{
// 				"2",
// 				"3",
// 				"1",
// 			},
// 			setup: func(myPoll *poll) {
// 				myPoll.ElectorIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 					"3",
// 				}
// 				myPoll.RemainingElectorIDs = []enum.PlayerID{
// 					"2",
// 					"3",
// 				}
// 			},
// 		},
// 		{
// 			name: "Ok (Doesn't exist in remaining and all)",
// 			electorIDs: []enum.PlayerID{
// 				"1",
// 				"2",
// 			},
// 			expectedStatus: true,
// 			newElectorIDs: []enum.PlayerID{
// 				"2",
// 				"3",
// 				"1",
// 			},
// 			newRemainingElectorIDs: []enum.PlayerID{
// 				"2",
// 				"3",
// 				"1",
// 			},
// 			setup: func(myPoll *poll) {
// 				myPoll.ElectorIDs = []enum.PlayerID{
// 					"2",
// 					"3",
// 				}
// 				myPoll.RemainingElectorIDs = []enum.PlayerID{
// 					"2",
// 					"3",
// 				}
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ps.Run(test.name, func() {
// 			myPoll, _ := NewPoll(capacity)
// 			test.setup(myPoll.(*poll))
// 			ok := myPoll.AddElectors(test.electorIDs...)

// 			ps.Equal(test.expectedStatus, ok)
// 			ps.Equal(test.newElectorIDs, myPoll.(*poll).ElectorIDs)
// 			ps.Equal(test.newRemainingElectorIDs, myPoll.(*poll).RemainingElectorIDs)
// 		})
// 	}
// }

// func (ps *PollSuite) TestRemoveElector() {
// 	tests := []struct {
// 		name                   string
// 		electorID              enum.PlayerID
// 		expectedStatus         bool
// 		newElectorIDs          []enum.PlayerID
// 		newRemainingElectorIDs []enum.PlayerID
// 		setup                  func(*poll)
// 	}{
// 		{
// 			name:           "Failure (Non-existent elector)",
// 			electorID:      "99",
// 			expectedStatus: false,
// 			newElectorIDs: []enum.PlayerID{
// 				"1",
// 				"2",
// 				"3",
// 			},
// 			newRemainingElectorIDs: []enum.PlayerID{
// 				"1",
// 				"2",
// 				"3",
// 			},
// 			setup: func(myPoll *poll) {
// 				myPoll.ElectorIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 					"3",
// 				}
// 				myPoll.RemainingElectorIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 					"3",
// 				}
// 			},
// 		},
// 		{
// 			name:           "Ok",
// 			electorID:      "1",
// 			expectedStatus: true,
// 			newElectorIDs: []enum.PlayerID{
// 				"1",
// 				"2",
// 				"3",
// 			},
// 			newRemainingElectorIDs: []enum.PlayerID{
// 				"2",
// 				"3",
// 			},
// 			setup: func(myPoll *poll) {
// 				myPoll.ElectorIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 					"3",
// 				}
// 				myPoll.RemainingElectorIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 					"3",
// 				}
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ps.Run(test.name, func() {
// 			myPoll, _ := NewPoll(9)
// 			test.setup(myPoll.(*poll))
// 			ok := myPoll.RemoveElector(test.electorID)

// 			ps.Equal(test.expectedStatus, ok)
// 			ps.Equal(test.newElectorIDs, myPoll.(*poll).ElectorIDs)
// 			ps.Equal(test.newRemainingElectorIDs, myPoll.(*poll).RemainingElectorIDs)
// 		})
// 	}
// }

// func (ps *PollSuite) TestSetWeight() {
// 	tests := []struct {
// 		name           string
// 		electorID      enum.PlayerID
// 		weight         uint
// 		expectedStatus bool
// 		setup          func(*poll)
// 	}{
// 		{
// 			name:           "Failure (Non-existent elector)",
// 			electorID:      "99",
// 			weight:         1,
// 			expectedStatus: false,
// 			setup: func(myPoll *poll) {
// 				myPoll.RemainingElectorIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 				}
// 			},
// 		},
// 		{
// 			name:           "Ok",
// 			electorID:      "1",
// 			weight:         5,
// 			expectedStatus: true,
// 			setup: func(myPoll *poll) {
// 				myPoll.RemainingElectorIDs = []enum.PlayerID{
// 					"1",
// 					"2",
// 				}
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ps.Run(test.name, func() {
// 			myPoll, _ := NewPoll(9)
// 			test.setup(myPoll.(*poll))
// 			ok := myPoll.SetWeight(test.electorID, test.weight)

// 			ps.Equal(test.expectedStatus, ok)

// 			if test.expectedStatus == true {
// 				ps.Equal(myPoll.(*poll).Weights[test.electorID], test.weight)
// 			}
// 		})
// 	}
// }

// func (ps *PollSuite) TestVote() {
// 	tests := []struct {
// 		name           string
// 		electorID      enum.PlayerID
// 		candidateID    enum.PlayerID
// 		expectedStatus bool
// 		expectedErr    string
// 		newWeight      uint
// 		newVotes       uint
// 		setup          func(*poll)
// 	}{
// 		{
// 			name:           "Failure (Cannot vote)",
// 			electorID:      "99",
// 			candidateID:    "2",
// 			expectedStatus: false,
// 			expectedErr:    "You're not allowed to vote ノ(ジ)ー'",
// 			setup: func(myPoll *poll) {
// 				// Open poll
// 				myPoll.Round = 1
// 				myPoll.Records[myPoll.Round] = &types.PollRecord{
// 					IsClosed: false,
// 				}

// 				myPoll.RemainingElectorIDs = []enum.PlayerID{"1"}
// 				myPoll.RemainingCandidateIDs = []enum.PlayerID{"2"}
// 			},
// 		},
// 		{
// 			name:           "Failure (Non-existent candidate)",
// 			electorID:      "1",
// 			candidateID:    "99",
// 			expectedStatus: false,
// 			expectedErr:    "Your vote is not valid ¬_¬",
// 			setup: func(myPoll *poll) {
// 				// Open poll
// 				myPoll.Round = 1
// 				myPoll.Records[myPoll.Round] = &types.PollRecord{
// 					IsClosed: false,
// 				}

// 				myPoll.RemainingElectorIDs = []enum.PlayerID{"1"}
// 				myPoll.RemainingCandidateIDs = []enum.PlayerID{"2"}
// 			},
// 		},
// 		{
// 			name:           "Ok (Skip)",
// 			electorID:      "1",
// 			candidateID:    "",
// 			expectedStatus: true,
// 			newVotes:       1,
// 			newWeight:      1,
// 			setup: func(myPoll *poll) {
// 				// Open poll
// 				myPoll.Round = 1
// 				myPoll.Records[myPoll.Round] = &types.PollRecord{
// 					IsClosed:    false,
// 					VoteRecords: make(map[enum.PlayerID]*types.VoteRecord),
// 				}

// 				myPoll.RemainingElectorIDs = []enum.PlayerID{"1"}
// 				myPoll.RemainingCandidateIDs = []enum.PlayerID{"2"}
// 				myPoll.SetWeight("1", 8)
// 			},
// 		},
// 		{
// 			name:           "Ok (Voted)",
// 			electorID:      "1",
// 			candidateID:    "2",
// 			expectedStatus: true,
// 			newVotes:       1,
// 			newWeight:      8,
// 			setup: func(myPoll *poll) {
// 				// Open poll
// 				myPoll.Round = 1
// 				myPoll.Records[myPoll.Round] = &types.PollRecord{
// 					IsClosed:    false,
// 					VoteRecords: make(map[enum.PlayerID]*types.VoteRecord),
// 				}

// 				myPoll.RemainingElectorIDs = []enum.PlayerID{"1"}
// 				myPoll.RemainingCandidateIDs = []enum.PlayerID{"2"}
// 				myPoll.SetWeight("1", 8)
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		ps.Run(test.name, func() {
// 			myPoll, _ := NewPoll(9)
// 			test.setup(myPoll.(*poll))
// 			ok, err := myPoll.Vote(test.electorID, test.candidateID)

// 			ps.Equal(test.expectedStatus, ok)

// 			if test.expectedStatus == true {
// 				ps.Equal(
// 					test.newVotes,
// 					myPoll.Record(enum.LastRound).VoteRecords[test.candidateID].Votes,
// 				)
// 				ps.Equal(
// 					test.newWeight,
// 					myPoll.Record(enum.LastRound).VoteRecords[test.candidateID].Weights,
// 				)
// 				ps.Contains(
// 					myPoll.Record(enum.LastRound).VoteRecords[test.candidateID].ElectorIDs,
// 					test.electorID,
// 				)
// 				ps.Contains(myPoll.(*poll).VotedElectorIDs, test.electorID)
// 			} else {
// 				ps.NotNil(err)
// 				ps.Equal(test.expectedErr, err.Error())
// 			}
// 		})
// 	}
// }
