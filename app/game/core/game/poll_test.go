package game

import (
	"testing"
	"uwwolf/app/game/config"
	"uwwolf/app/game/types"

	"github.com/stretchr/testify/assert"
)

func TestNewPoll(t *testing.T) {
	tests := []struct {
		name             string
		expectedCapacity uint
		expectedErr      string
	}{
		{
			name:             "Too small capacity",
			expectedCapacity: 2,
			expectedErr:      "The capacity is too small ¬_¬",
		},
		{
			name:             "Ok",
			expectedCapacity: 5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			myPoll, err := NewPoll(test.expectedCapacity)

			if test.expectedErr != "" {
				assert.Equal(t, test.expectedErr, err.Error())
			} else {
				assert.NotNil(t, myPoll)
				assert.NotNil(t, myPoll.(*poll).Weights)
				assert.NotNil(t, myPoll.(*poll).Records)
				assert.Equal(t, test.expectedCapacity, myPoll.(*poll).Capacity)
			}
		})
	}
}

func TestIsOpenPoll(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus bool
		expectedRound  types.Round
		setup          func(*poll)
	}{
		{
			name:           "Not open yet (Round is zero)",
			expectedStatus: false,
			expectedRound:  0,
			setup: func(myPoll *poll) {
				myPoll.Round = 0
			},
		},
		{
			name:           "Not open yet (Poll record is nil)",
			expectedStatus: false,
			expectedRound:  1,
			setup: func(myPoll *poll) {
				myPoll.Round = 1
			},
		},
		{
			name:           "Not open yet (Poll was closed)",
			expectedStatus: false,
			expectedRound:  1,
			setup: func(myPoll *poll) {
				myPoll.Round = 1
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed: true,
				}
			},
		},
		{
			name:           "Is open",
			expectedStatus: true,
			expectedRound:  1,
			setup: func(myPoll *poll) {
				myPoll.Round = 1
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed: false,
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			myPoll, _ := NewPoll(9)
			test.setup(myPoll.(*poll))
			isOpen, round := myPoll.IsOpen()

			assert.Equal(t, test.expectedStatus, isOpen)
			assert.Equal(t, test.expectedRound, round)
		})
	}
}

func TestCanVotePoll(t *testing.T) {
	electorID := types.PlayerID("1")
	tests := []struct {
		name        string
		electorID   types.PlayerID
		expectedRes bool
		setup       func(*poll)
	}{
		{
			name:        "Cannot vote (Poll was closed)",
			electorID:   electorID,
			expectedRes: false,
			setup: func(myPoll *poll) {
				myPoll.Round = 0
			},
		},
		{
			name:        "Cannot vote (Not an elector)",
			electorID:   electorID,
			expectedRes: false,
			setup: func(myPoll *poll) {
				// Open poll
				myPoll.Round = 1
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed: false,
				}
			},
		},
		{
			name:        "Cannot vote (Already voted)",
			electorID:   electorID,
			expectedRes: false,
			setup: func(myPoll *poll) {
				// Open poll
				myPoll.Round = 1
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed: false,
				}

				myPoll.VotedElectorIDs = append(
					myPoll.VotedElectorIDs,
					electorID,
				)
			},
		},
		{
			name:        "Ok",
			electorID:   electorID,
			expectedRes: true,
			setup: func(myPoll *poll) {
				// Open poll
				myPoll.Round = 1
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed: false,
				}

				myPoll.RemainingElectorIDs = append(
					myPoll.RemainingElectorIDs,
					electorID,
				)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			myPoll, _ := NewPoll(9)
			test.setup(myPoll.(*poll))
			res := myPoll.CanVote(test.electorID)

			assert.Equal(t, test.expectedRes, res)
		})
	}
}

func TestRecordPoll(t *testing.T) {
	tests := []struct {
		name           string
		round          types.Round
		expecredRecord *types.PollRecord
		setup          func(*poll)
	}{
		{
			name:           "Nil (Round is zero)",
			round:          0,
			expecredRecord: nil,
			setup:          func(myPoll *poll) {},
		},
		{
			name:           "Nil (Non-existent round)",
			round:          99,
			expecredRecord: nil,
			setup: func(myPoll *poll) {
				myPoll.Round = 1
				myPoll.Records[myPoll.Round] = &types.PollRecord{}
			},
		},
		{
			name:  "Ok",
			round: 1,
			expecredRecord: &types.PollRecord{
				WinnerID: "98",
			},
			setup: func(myPoll *poll) {
				myPoll.Round = 1
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					WinnerID: "98",
				}
			},
		},
		{
			name:  "Ok (Latest record)",
			round: config.LastRound,
			expecredRecord: &types.PollRecord{
				WinnerID: "99",
			},
			setup: func(myPoll *poll) {
				myPoll.Records[1] = &types.PollRecord{
					WinnerID: "98",
				}
				myPoll.Round = 2
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					WinnerID: "99",
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			myPoll, _ := NewPoll(9)
			test.setup(myPoll.(*poll))
			record := myPoll.Record(test.round)

			assert.Equal(t, test.expecredRecord, record)
		})
	}
}

func TestOpenPoll(t *testing.T) {
	capacity := uint(3)
	tests := []struct {
		name           string
		expectedStatus bool
		expectedRound  types.Round
		setup          func(*poll)
	}{
		{
			name:           "Cannot open (Already open)",
			expectedStatus: false,
			expectedRound:  1,
			setup: func(myPoll *poll) {
				myPoll.Round = 1
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed: false,
				}
			},
		},
		{
			name:           "Cannot open (Not enough electors)",
			expectedStatus: false,
			expectedRound:  1,
			setup: func(myPoll *poll) {
				myPoll.Round = 1
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed: true,
				}
				myPoll.RemainingElectorIDs = []types.PlayerID{}
			},
		},
		{
			name:           "Ok",
			expectedStatus: true,
			expectedRound:  2,
			setup: func(myPoll *poll) {
				myPoll.Round = 1
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed: true,
				}
				myPoll.RemainingElectorIDs = []types.PlayerID{"1", "2", "3"}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			myPoll, _ := NewPoll(capacity)
			test.setup(myPoll.(*poll))
			ok, round := myPoll.Open()

			assert.Equal(t, test.expectedStatus, ok)
			assert.Equal(t, test.expectedRound, round)

			if test.expectedStatus == true {
				assert.Empty(t, myPoll.(*poll).VotedElectorIDs)
				assert.False(t, myPoll.Record(config.LastRound).IsClosed)
			}
		})
	}
}

func TestClosePoll(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus bool
		expectedRecord *types.PollRecord
		setup          func(*poll)
	}{
		{
			name:           "Poll was closed",
			expectedStatus: false,
			expectedRecord: nil,
			setup: func(myPoll *poll) {
				myPoll.Round = 1
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed: true,
				}
			},
		},
		{
			name:           "Ok with a small weight for candidates",
			expectedStatus: true,
			expectedRecord: &types.PollRecord{
				IsClosed: true,
				WinnerID: "",
				VoteRecords: map[types.PlayerID]*types.VoteRecord{
					"1": {
						ElectorIDs: []types.PlayerID{"2", "3"},
						Weights:    2,
						Votes:      2,
					},
					"2": {
						ElectorIDs: []types.PlayerID{"1", "4"},
						Weights:    2,
						Votes:      2,
					},
					"": {
						ElectorIDs: []types.PlayerID{"5"},
						Votes:      1,
						Weights:    1,
					},
				},
			},
			setup: func(myPoll *poll) {
				myPoll.Round = 1
				myPoll.RemainingElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
					"4",
					"5",
				}
				myPoll.VotedElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
					"4",
				}
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed: false,
					VoteRecords: map[types.PlayerID]*types.VoteRecord{
						"1": {
							ElectorIDs: []types.PlayerID{"2", "3"},
							Weights:    2,
							Votes:      2,
						},
						"2": {
							ElectorIDs: []types.PlayerID{"1", "4"},
							Weights:    2,
							Votes:      2,
						},
					},
				}
				myPoll.Records[myPoll.Round].VoteRecords[""] = &types.VoteRecord{
					ElectorIDs: []types.PlayerID{},
				}
			},
		},
		{
			name:           "Ok with a draw",
			expectedStatus: true,
			expectedRecord: &types.PollRecord{
				IsClosed: true,
				WinnerID: "",
				VoteRecords: map[types.PlayerID]*types.VoteRecord{
					"1": {
						ElectorIDs: []types.PlayerID{"2", "3", "4"},
						Weights:    3,
						Votes:      3,
					},
					"2": {
						ElectorIDs: []types.PlayerID{"1", "5", "6"},
						Weights:    3,
						Votes:      3,
					},
					"": {
						ElectorIDs: []types.PlayerID{},
						Votes:      0,
						Weights:    0,
					},
				},
			},
			setup: func(myPoll *poll) {
				myPoll.Round = 1
				myPoll.RemainingElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
					"4",
					"5",
					"6",
				}
				myPoll.VotedElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
					"4",
					"5",
					"6",
				}
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed: false,
					VoteRecords: map[types.PlayerID]*types.VoteRecord{
						"1": {
							ElectorIDs: []types.PlayerID{"2", "3", "4"},
							Weights:    3,
							Votes:      3,
						},
						"2": {
							ElectorIDs: []types.PlayerID{"1", "5", "6"},
							Weights:    3,
							Votes:      3,
						},
					},
				}
				myPoll.Records[myPoll.Round].VoteRecords[""] = &types.VoteRecord{
					ElectorIDs: []types.PlayerID{},
				}
			},
		},
		{
			name:           "Ok with a winner",
			expectedStatus: true,
			expectedRecord: &types.PollRecord{
				IsClosed: true,
				WinnerID: "1",
				VoteRecords: map[types.PlayerID]*types.VoteRecord{
					"1": {
						ElectorIDs: []types.PlayerID{"2", "3", "4"},
						Weights:    3,
						Votes:      3,
					},
					"2": {
						ElectorIDs: []types.PlayerID{"1"},
						Weights:    1,
						Votes:      1,
					},
					"": {
						ElectorIDs: []types.PlayerID{"5"},
						Votes:      1,
						Weights:    1,
					},
				},
			},
			setup: func(myPoll *poll) {
				myPoll.Round = 1
				myPoll.RemainingElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
					"4",
					"5",
				}
				myPoll.VotedElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
					"4",
				}
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed: false,
					VoteRecords: map[types.PlayerID]*types.VoteRecord{
						"1": {
							ElectorIDs: []types.PlayerID{"2", "3", "4"},
							Weights:    3,
							Votes:      3,
						},
						"2": {
							ElectorIDs: []types.PlayerID{"1"},
							Weights:    1,
							Votes:      1,
						},
					},
				}
				myPoll.Records[myPoll.Round].VoteRecords[""] = &types.VoteRecord{
					ElectorIDs: []types.PlayerID{},
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			myPoll, _ := NewPoll(9)
			test.setup(myPoll.(*poll))
			ok, record := myPoll.Close()

			assert.Equal(t, test.expectedStatus, ok)
			assert.Equal(t, test.expectedRecord, record)

			if test.expectedStatus == true {
				assert.Equal(
					t,
					len(myPoll.(*poll).RemainingElectorIDs),
					len(myPoll.(*poll).VotedElectorIDs),
				)

				isOpen, _ := myPoll.IsOpen()
				assert.False(t, isOpen)
			}
		})
	}
}

func TestAddCandidatesPoll(t *testing.T) {
	tests := []struct {
		name                          string
		candidateIDs                  []types.PlayerID
		expectedCandidateIDs          []types.PlayerID
		expectedRemainingCandidateIDs []types.PlayerID
		setup                         func(*poll)
	}{
		{
			name: "Exist in remaining",
			candidateIDs: []types.PlayerID{
				"1",
				"2",
			},
			expectedCandidateIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			expectedRemainingCandidateIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			setup: func(myPoll *poll) {
				myPoll.CandidateIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
				myPoll.RemainingCandidateIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
			},
		},
		{
			name: "Doesn't exist in remaining but in all",
			candidateIDs: []types.PlayerID{
				"1",
				"2",
			},
			expectedCandidateIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			expectedRemainingCandidateIDs: []types.PlayerID{
				"2",
				"3",
				"1",
			},
			setup: func(myPoll *poll) {
				myPoll.CandidateIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
				myPoll.RemainingCandidateIDs = []types.PlayerID{
					"2",
					"3",
				}
			},
		},
		{
			name: "Doesn't exist in remaining and all",
			candidateIDs: []types.PlayerID{
				"1",
				"2",
			},
			expectedCandidateIDs: []types.PlayerID{
				"2",
				"3",
				"1",
			},
			expectedRemainingCandidateIDs: []types.PlayerID{
				"2",
				"3",
				"1",
			},
			setup: func(myPoll *poll) {
				myPoll.CandidateIDs = []types.PlayerID{
					"2",
					"3",
				}
				myPoll.RemainingCandidateIDs = []types.PlayerID{
					"2",
					"3",
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			myPoll, _ := NewPoll(9)
			test.setup(myPoll.(*poll))
			myPoll.AddCandidates(test.candidateIDs...)

			assert.Equal(t, test.expectedCandidateIDs, myPoll.(*poll).CandidateIDs)
			assert.Equal(t, test.expectedRemainingCandidateIDs, myPoll.(*poll).RemainingCandidateIDs)
		})
	}
}

func TestRemoveCandidatePoll(t *testing.T) {
	tests := []struct {
		name                          string
		candidateID                   types.PlayerID
		expectedStatus                bool
		expectedCandidateIDs          []types.PlayerID
		expectedRemainingCandidateIDs []types.PlayerID
		setup                         func(*poll)
	}{
		{
			name:           "Remove non-existent candidate",
			candidateID:    "99",
			expectedStatus: false,
			expectedCandidateIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			expectedRemainingCandidateIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			setup: func(myPoll *poll) {
				myPoll.CandidateIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
				myPoll.RemainingCandidateIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
			},
		},
		{
			name:           "Ok",
			candidateID:    "1",
			expectedStatus: true,
			expectedCandidateIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			expectedRemainingCandidateIDs: []types.PlayerID{
				"2",
				"3",
			},
			setup: func(myPoll *poll) {
				myPoll.CandidateIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
				myPoll.RemainingCandidateIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			myPoll, _ := NewPoll(9)
			test.setup(myPoll.(*poll))
			ok := myPoll.RemoveCandidate(test.candidateID)

			assert.Equal(t, test.expectedStatus, ok)
			assert.Equal(t, test.expectedCandidateIDs, myPoll.(*poll).CandidateIDs)
			assert.Equal(t, test.expectedRemainingCandidateIDs, myPoll.(*poll).RemainingCandidateIDs)
		})
	}
}

func TestAddElectorsPoll(t *testing.T) {
	capacity := uint(5)
	tests := []struct {
		name                        string
		electorIDs                  []types.PlayerID
		expectedStatus              bool
		expectedElectorIDs          []types.PlayerID
		expectedRemainingElectorIDs []types.PlayerID
		setup                       func(*poll)
	}{
		{
			name: "Overload",
			electorIDs: []types.PlayerID{
				"4",
				"5",
				"7",
				"8",
			},
			expectedStatus: false,
			expectedElectorIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			expectedRemainingElectorIDs: []types.PlayerID{
				"1",
				"2",
			},
			setup: func(myPoll *poll) {
				myPoll.ElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
				myPoll.RemainingElectorIDs = []types.PlayerID{
					"1",
					"2",
				}
			},
		},
		{
			name: "Exist in remaining",
			electorIDs: []types.PlayerID{
				"1",
				"2",
			},
			expectedStatus: true,
			expectedElectorIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			expectedRemainingElectorIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			setup: func(myPoll *poll) {
				myPoll.ElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
				myPoll.RemainingElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
			},
		},
		{
			name: "Doesn't exist in remaining but in all",
			electorIDs: []types.PlayerID{
				"1",
				"2",
			},
			expectedStatus: true,
			expectedElectorIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			expectedRemainingElectorIDs: []types.PlayerID{
				"2",
				"3",
				"1",
			},
			setup: func(myPoll *poll) {
				myPoll.ElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
				myPoll.RemainingElectorIDs = []types.PlayerID{
					"2",
					"3",
				}
			},
		},
		{
			name: "Doesn't exist in remaining and all",
			electorIDs: []types.PlayerID{
				"1",
				"2",
			},
			expectedStatus: true,
			expectedElectorIDs: []types.PlayerID{
				"2",
				"3",
				"1",
			},
			expectedRemainingElectorIDs: []types.PlayerID{
				"2",
				"3",
				"1",
			},
			setup: func(myPoll *poll) {
				myPoll.ElectorIDs = []types.PlayerID{
					"2",
					"3",
				}
				myPoll.RemainingElectorIDs = []types.PlayerID{
					"2",
					"3",
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			myPoll, _ := NewPoll(capacity)
			test.setup(myPoll.(*poll))
			ok := myPoll.AddElectors(test.electorIDs...)

			assert.Equal(t, test.expectedStatus, ok)
			assert.Equal(t, test.expectedElectorIDs, myPoll.(*poll).ElectorIDs)
			assert.Equal(t, test.expectedRemainingElectorIDs, myPoll.(*poll).RemainingElectorIDs)
		})
	}
}

func TestRemoveElectorPoll(t *testing.T) {
	tests := []struct {
		name                        string
		electorID                   types.PlayerID
		expectedStatus              bool
		expectedElectorIDs          []types.PlayerID
		expectedRemainingElectorIDs []types.PlayerID
		setup                       func(*poll)
	}{
		{
			name:           "Remove non-existent elector",
			electorID:      "99",
			expectedStatus: false,
			expectedElectorIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			expectedRemainingElectorIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			setup: func(myPoll *poll) {
				myPoll.ElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
				myPoll.RemainingElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
			},
		},
		{
			name:           "Ok",
			electorID:      "1",
			expectedStatus: true,
			expectedElectorIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			expectedRemainingElectorIDs: []types.PlayerID{
				"2",
				"3",
			},
			setup: func(myPoll *poll) {
				myPoll.ElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
				myPoll.RemainingElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			myPoll, _ := NewPoll(9)
			test.setup(myPoll.(*poll))
			ok := myPoll.RemoveElector(test.electorID)

			assert.Equal(t, test.expectedStatus, ok)
			assert.Equal(t, test.expectedElectorIDs, myPoll.(*poll).ElectorIDs)
			assert.Equal(t, test.expectedRemainingElectorIDs, myPoll.(*poll).RemainingElectorIDs)
		})
	}
}

func TestSetWeightPoll(t *testing.T) {
	tests := []struct {
		name           string
		electorID      types.PlayerID
		weight         uint
		expectedStatus bool
		setup          func(*poll)
	}{
		{
			name:           "Non-existent elector",
			electorID:      "99",
			weight:         1,
			expectedStatus: false,
			setup: func(myPoll *poll) {
				myPoll.RemainingElectorIDs = []types.PlayerID{
					"1",
					"2",
				}
			},
		},
		{
			name:           "Ok",
			electorID:      "1",
			weight:         5,
			expectedStatus: true,
			setup: func(myPoll *poll) {
				myPoll.RemainingElectorIDs = []types.PlayerID{
					"1",
					"2",
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			myPoll, _ := NewPoll(9)
			test.setup(myPoll.(*poll))
			ok := myPoll.SetWeight(test.electorID, test.weight)

			assert.Equal(t, test.expectedStatus, ok)

			if test.expectedStatus == true {
				assert.Equal(t, myPoll.(*poll).Weights[test.electorID], test.weight)
			}
		})
	}
}

func TestVotePoll(t *testing.T) {
	tests := []struct {
		name           string
		electorID      types.PlayerID
		candidateID    types.PlayerID
		expectedStatus bool
		expectedWeight uint
		expectedVotes  uint
		setup          func(*poll)
	}{
		{
			name:           "Cannot vote",
			electorID:      "99",
			candidateID:    "2",
			expectedStatus: false,
			setup: func(myPoll *poll) {
				// Open poll
				myPoll.Round = 1
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed: false,
				}

				myPoll.RemainingElectorIDs = []types.PlayerID{"1"}
				myPoll.RemainingCandidateIDs = []types.PlayerID{"2"}
			},
		},
		{
			name:           "Non-existent candidate",
			electorID:      "1",
			candidateID:    "99",
			expectedStatus: false,
			setup: func(myPoll *poll) {
				// Open poll
				myPoll.Round = 1
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed: false,
				}

				myPoll.RemainingElectorIDs = []types.PlayerID{"1"}
				myPoll.RemainingCandidateIDs = []types.PlayerID{"2"}
			},
		},
		{
			name:           "Ok (Skip)",
			electorID:      "1",
			candidateID:    "",
			expectedStatus: true,
			expectedVotes:  1,
			expectedWeight: 1,
			setup: func(myPoll *poll) {
				// Open poll
				myPoll.Round = 1
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed:    false,
					VoteRecords: make(map[types.PlayerID]*types.VoteRecord),
				}

				myPoll.RemainingElectorIDs = []types.PlayerID{"1"}
				myPoll.RemainingCandidateIDs = []types.PlayerID{"2"}
				myPoll.SetWeight("1", 8)
			},
		},
		{
			name:           "Ok (Voted)",
			electorID:      "1",
			candidateID:    "2",
			expectedStatus: true,
			expectedVotes:  1,
			expectedWeight: 8,
			setup: func(myPoll *poll) {
				// Open poll
				myPoll.Round = 1
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed:    false,
					VoteRecords: make(map[types.PlayerID]*types.VoteRecord),
				}

				myPoll.RemainingElectorIDs = []types.PlayerID{"1"}
				myPoll.RemainingCandidateIDs = []types.PlayerID{"2"}
				myPoll.SetWeight("1", 8)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			myPoll, _ := NewPoll(9)
			test.setup(myPoll.(*poll))
			ok := myPoll.Vote(test.electorID, test.candidateID)

			assert.Equal(t, test.expectedStatus, ok)

			if test.expectedStatus == true {
				assert.Equal(t,
					test.expectedVotes,
					myPoll.Record(config.LastRound).VoteRecords[test.candidateID].Votes,
				)
				assert.Equal(t,
					test.expectedWeight,
					myPoll.Record(config.LastRound).VoteRecords[test.candidateID].Weights,
				)
				assert.Contains(
					t,
					myPoll.Record(config.LastRound).VoteRecords[test.candidateID].ElectorIDs,
					test.electorID,
				)
				assert.Contains(t, myPoll.(*poll).VotedElectorIDs, test.electorID)
			}
		})
	}
}
