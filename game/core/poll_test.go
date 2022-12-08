package core

import (
	"testing"
	"uwwolf/game/enum"
	"uwwolf/game/types"

	"github.com/stretchr/testify/assert"
)

func TestNewPoll(t *testing.T) {
	tests := []struct {
		name        string
		capacity    uint8
		expectedErr string
	}{
		{
			name:        "Failure (Too small capacity)",
			capacity:    2,
			expectedErr: "The capacity is too small ¬_¬",
		},
		{
			name:     "Ok",
			capacity: 5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			myPoll, err := NewPoll(test.capacity)

			if test.expectedErr != "" {
				assert.Equal(t, test.expectedErr, err.Error())
			} else {
				assert.NotNil(t, myPoll)
				assert.NotNil(t, myPoll.(*poll).Weights)
				assert.NotNil(t, myPoll.(*poll).Records)
				assert.Equal(t, test.capacity, myPoll.(*poll).Capacity)
			}
		})
	}
}

func TestIsOpenPoll(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus bool
		expectedRound  enum.Round
		setup          func(*poll)
	}{
		{
			name:           "Not open (Round is zero)",
			expectedStatus: false,
			expectedRound:  0,
			setup: func(myPoll *poll) {
				myPoll.Round = 0
			},
		},
		{
			name:           "Not open (Poll record is nil)",
			expectedStatus: false,
			expectedRound:  1,
			setup: func(myPoll *poll) {
				myPoll.Round = 1
			},
		},
		{
			name:           "Not open (Poll was closed)",
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
	electorID := enum.PlayerID("1")
	tests := []struct {
		name        string
		electorID   enum.PlayerID
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
		round          enum.Round
		expectedRecord *types.PollRecord
		setup          func(*poll)
	}{
		{
			name:           "Nil (Round is zero)",
			round:          0,
			expectedRecord: nil,
			setup:          func(myPoll *poll) {},
		},
		{
			name:           "Nil (Non-existent round)",
			round:          99,
			expectedRecord: nil,
			setup: func(myPoll *poll) {
				myPoll.Round = 1
				myPoll.Records[myPoll.Round] = &types.PollRecord{}
			},
		},
		{
			name:  "Ok",
			round: 1,
			expectedRecord: &types.PollRecord{
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
			name:  "Ok (Get latest record)",
			round: enum.LastRound,
			expectedRecord: &types.PollRecord{
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

			assert.Equal(t, test.expectedRecord, record)
		})
	}
}

func TestOpenPoll(t *testing.T) {
	capacity := uint8(5)
	tests := []struct {
		name           string
		expectedStatus bool
		expectedRound  enum.Round
		setup          func(*poll)
	}{
		{
			name:           "Failure (Already open)",
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
			name:           "Failure (Not enough electors)",
			expectedStatus: false,
			expectedRound:  1,
			setup: func(myPoll *poll) {
				myPoll.Round = 1
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed: true,
				}
				myPoll.RemainingElectorIDs = []enum.PlayerID{"1", "2", "3"}
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
				myPoll.VotedElectorIDs = []enum.PlayerID{"1", "2", "3", "4", "5"}
				myPoll.RemainingElectorIDs = []enum.PlayerID{"1", "2", "3", "4", "5"}
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
				assert.False(t, myPoll.Record(enum.LastRound).IsClosed)
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
			name:           "Failure (Poll was closed)",
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
			name:           "Ok (No one have half-votes weight)",
			expectedStatus: true,
			expectedRecord: &types.PollRecord{
				IsClosed: true,
				WinnerID: "",
				VoteRecords: map[enum.PlayerID]*types.VoteRecord{
					"1": {
						ElectorIDs: []enum.PlayerID{"2", "3"},
						Weights:    2,
						Votes:      2,
					},
					"2": {
						ElectorIDs: []enum.PlayerID{"1", "4"},
						Weights:    2,
						Votes:      2,
					},
					"": {
						ElectorIDs: []enum.PlayerID{"5"},
						Votes:      1,
						Weights:    1,
					},
				},
			},
			setup: func(myPoll *poll) {
				myPoll.Round = 1
				myPoll.RemainingElectorIDs = []enum.PlayerID{
					"1",
					"2",
					"3",
					"4",
					"5",
				}
				myPoll.VotedElectorIDs = []enum.PlayerID{
					"1",
					"2",
					"3",
					"4",
				}
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed: false,
					VoteRecords: map[enum.PlayerID]*types.VoteRecord{
						"1": {
							ElectorIDs: []enum.PlayerID{"2", "3"},
							Weights:    2,
							Votes:      2,
						},
						"2": {
							ElectorIDs: []enum.PlayerID{"1", "4"},
							Weights:    2,
							Votes:      2,
						},
					},
				}
				myPoll.Records[myPoll.Round].VoteRecords[""] = &types.VoteRecord{
					ElectorIDs: []enum.PlayerID{},
				}
			},
		},
		{
			name:           "Ok (A draw)",
			expectedStatus: true,
			expectedRecord: &types.PollRecord{
				IsClosed: true,
				WinnerID: "",
				VoteRecords: map[enum.PlayerID]*types.VoteRecord{
					"1": {
						ElectorIDs: []enum.PlayerID{"2", "3", "4"},
						Weights:    3,
						Votes:      3,
					},
					"2": {
						ElectorIDs: []enum.PlayerID{"1", "5", "6"},
						Weights:    3,
						Votes:      3,
					},
					"": {
						ElectorIDs: []enum.PlayerID{},
						Votes:      0,
						Weights:    0,
					},
				},
			},
			setup: func(myPoll *poll) {
				myPoll.Round = 1
				myPoll.RemainingElectorIDs = []enum.PlayerID{
					"1",
					"2",
					"3",
					"4",
					"5",
					"6",
				}
				myPoll.VotedElectorIDs = []enum.PlayerID{
					"1",
					"2",
					"3",
					"4",
					"5",
					"6",
				}
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed: false,
					VoteRecords: map[enum.PlayerID]*types.VoteRecord{
						"1": {
							ElectorIDs: []enum.PlayerID{"2", "3", "4"},
							Weights:    3,
							Votes:      3,
						},
						"2": {
							ElectorIDs: []enum.PlayerID{"1", "5", "6"},
							Weights:    3,
							Votes:      3,
						},
					},
				}
				myPoll.Records[myPoll.Round].VoteRecords[""] = &types.VoteRecord{
					ElectorIDs: []enum.PlayerID{},
				}
			},
		},
		{
			name:           "Ok (One cadidate have half-votes weight)",
			expectedStatus: true,
			expectedRecord: &types.PollRecord{
				IsClosed: true,
				WinnerID: "1",
				VoteRecords: map[enum.PlayerID]*types.VoteRecord{
					"1": {
						ElectorIDs: []enum.PlayerID{"2", "3", "4"},
						Weights:    3,
						Votes:      3,
					},
					"2": {
						ElectorIDs: []enum.PlayerID{"1"},
						Weights:    1,
						Votes:      1,
					},
					"": {
						ElectorIDs: []enum.PlayerID{"5"},
						Votes:      1,
						Weights:    1,
					},
				},
			},
			setup: func(myPoll *poll) {
				myPoll.Round = 1
				myPoll.RemainingElectorIDs = []enum.PlayerID{
					"1",
					"2",
					"3",
					"4",
					"5",
				}
				myPoll.VotedElectorIDs = []enum.PlayerID{
					"1",
					"2",
					"3",
					"4",
				}
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed: false,
					VoteRecords: map[enum.PlayerID]*types.VoteRecord{
						"1": {
							ElectorIDs: []enum.PlayerID{"2", "3", "4"},
							Weights:    3,
							Votes:      3,
						},
						"2": {
							ElectorIDs: []enum.PlayerID{"1"},
							Weights:    1,
							Votes:      1,
						},
					},
				}
				myPoll.Records[myPoll.Round].VoteRecords[""] = &types.VoteRecord{
					ElectorIDs: []enum.PlayerID{},
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
		name                     string
		candidateIDs             []enum.PlayerID
		newCandidateIDs          []enum.PlayerID
		newRemainingCandidateIDs []enum.PlayerID
		setup                    func(*poll)
	}{
		{
			name: "Failure (Already existed in remaining)",
			candidateIDs: []enum.PlayerID{
				"1",
				"2",
			},
			newCandidateIDs: []enum.PlayerID{
				"1",
				"2",
				"3",
			},
			newRemainingCandidateIDs: []enum.PlayerID{
				"1",
				"2",
				"3",
			},
			setup: func(myPoll *poll) {
				myPoll.CandidateIDs = []enum.PlayerID{
					"1",
					"2",
					"3",
				}
				myPoll.RemainingCandidateIDs = []enum.PlayerID{
					"1",
					"2",
					"3",
				}
			},
		},
		{
			name: "Ok (Doesn't exist in remaining but in all)",
			candidateIDs: []enum.PlayerID{
				"1",
				"2",
			},
			newCandidateIDs: []enum.PlayerID{
				"1",
				"2",
				"3",
			},
			newRemainingCandidateIDs: []enum.PlayerID{
				"2",
				"3",
				"1",
			},
			setup: func(myPoll *poll) {
				myPoll.CandidateIDs = []enum.PlayerID{
					"1",
					"2",
					"3",
				}
				myPoll.RemainingCandidateIDs = []enum.PlayerID{
					"2",
					"3",
				}
			},
		},
		{
			name: "Ok (Doesn't exist in remaining and all)",
			candidateIDs: []enum.PlayerID{
				"1",
				"2",
			},
			newCandidateIDs: []enum.PlayerID{
				"2",
				"3",
				"1",
			},
			newRemainingCandidateIDs: []enum.PlayerID{
				"2",
				"3",
				"1",
			},
			setup: func(myPoll *poll) {
				myPoll.CandidateIDs = []enum.PlayerID{
					"2",
					"3",
				}
				myPoll.RemainingCandidateIDs = []enum.PlayerID{
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

			assert.Equal(t, test.newCandidateIDs, myPoll.(*poll).CandidateIDs)
			assert.Equal(t, test.newRemainingCandidateIDs, myPoll.(*poll).RemainingCandidateIDs)
		})
	}
}

func TestRemoveCandidatePoll(t *testing.T) {
	tests := []struct {
		name                     string
		candidateID              enum.PlayerID
		expectedStatus           bool
		newCandidateIDs          []enum.PlayerID
		newRemainingCandidateIDs []enum.PlayerID
		setup                    func(*poll)
	}{
		{
			name:           "Failure (Non-existent candidate)",
			candidateID:    "99",
			expectedStatus: false,
			newCandidateIDs: []enum.PlayerID{
				"1",
				"2",
				"3",
			},
			newRemainingCandidateIDs: []enum.PlayerID{
				"1",
				"2",
				"3",
			},
			setup: func(myPoll *poll) {
				myPoll.CandidateIDs = []enum.PlayerID{
					"1",
					"2",
					"3",
				}
				myPoll.RemainingCandidateIDs = []enum.PlayerID{
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
			newCandidateIDs: []enum.PlayerID{
				"1",
				"2",
				"3",
			},
			newRemainingCandidateIDs: []enum.PlayerID{
				"2",
				"3",
			},
			setup: func(myPoll *poll) {
				myPoll.CandidateIDs = []enum.PlayerID{
					"1",
					"2",
					"3",
				}
				myPoll.RemainingCandidateIDs = []enum.PlayerID{
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
			assert.Equal(t, test.newCandidateIDs, myPoll.(*poll).CandidateIDs)
			assert.Equal(t, test.newRemainingCandidateIDs, myPoll.(*poll).RemainingCandidateIDs)
		})
	}
}

func TestAddElectorsPoll(t *testing.T) {
	capacity := uint8(5)
	tests := []struct {
		name                   string
		electorIDs             []enum.PlayerID
		expectedStatus         bool
		newElectorIDs          []enum.PlayerID
		newRemainingElectorIDs []enum.PlayerID
		setup                  func(*poll)
	}{
		{
			name: "Failure (Overload)",
			electorIDs: []enum.PlayerID{
				"4",
				"5",
				"7",
				"8",
			},
			expectedStatus: false,
			newElectorIDs: []enum.PlayerID{
				"1",
				"2",
				"3",
			},
			newRemainingElectorIDs: []enum.PlayerID{
				"1",
				"2",
			},
			setup: func(myPoll *poll) {
				myPoll.ElectorIDs = []enum.PlayerID{
					"1",
					"2",
					"3",
				}
				myPoll.RemainingElectorIDs = []enum.PlayerID{
					"1",
					"2",
				}
			},
		},
		{
			name: "Failure (Already existed in remaining)",
			electorIDs: []enum.PlayerID{
				"1",
				"2",
			},
			expectedStatus: true,
			newElectorIDs: []enum.PlayerID{
				"1",
				"2",
				"3",
			},
			newRemainingElectorIDs: []enum.PlayerID{
				"1",
				"2",
				"3",
			},
			setup: func(myPoll *poll) {
				myPoll.ElectorIDs = []enum.PlayerID{
					"1",
					"2",
					"3",
				}
				myPoll.RemainingElectorIDs = []enum.PlayerID{
					"1",
					"2",
					"3",
				}
			},
		},
		{
			name: "Ok (Doesn't exist in remaining but in all)",
			electorIDs: []enum.PlayerID{
				"1",
				"2",
			},
			expectedStatus: true,
			newElectorIDs: []enum.PlayerID{
				"1",
				"2",
				"3",
			},
			newRemainingElectorIDs: []enum.PlayerID{
				"2",
				"3",
				"1",
			},
			setup: func(myPoll *poll) {
				myPoll.ElectorIDs = []enum.PlayerID{
					"1",
					"2",
					"3",
				}
				myPoll.RemainingElectorIDs = []enum.PlayerID{
					"2",
					"3",
				}
			},
		},
		{
			name: "Ok (Doesn't exist in remaining and all)",
			electorIDs: []enum.PlayerID{
				"1",
				"2",
			},
			expectedStatus: true,
			newElectorIDs: []enum.PlayerID{
				"2",
				"3",
				"1",
			},
			newRemainingElectorIDs: []enum.PlayerID{
				"2",
				"3",
				"1",
			},
			setup: func(myPoll *poll) {
				myPoll.ElectorIDs = []enum.PlayerID{
					"2",
					"3",
				}
				myPoll.RemainingElectorIDs = []enum.PlayerID{
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
			assert.Equal(t, test.newElectorIDs, myPoll.(*poll).ElectorIDs)
			assert.Equal(t, test.newRemainingElectorIDs, myPoll.(*poll).RemainingElectorIDs)
		})
	}
}

func TestRemoveElectorPoll(t *testing.T) {
	tests := []struct {
		name                   string
		electorID              enum.PlayerID
		expectedStatus         bool
		newElectorIDs          []enum.PlayerID
		newRemainingElectorIDs []enum.PlayerID
		setup                  func(*poll)
	}{
		{
			name:           "Failure (Non-existent elector)",
			electorID:      "99",
			expectedStatus: false,
			newElectorIDs: []enum.PlayerID{
				"1",
				"2",
				"3",
			},
			newRemainingElectorIDs: []enum.PlayerID{
				"1",
				"2",
				"3",
			},
			setup: func(myPoll *poll) {
				myPoll.ElectorIDs = []enum.PlayerID{
					"1",
					"2",
					"3",
				}
				myPoll.RemainingElectorIDs = []enum.PlayerID{
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
			newElectorIDs: []enum.PlayerID{
				"1",
				"2",
				"3",
			},
			newRemainingElectorIDs: []enum.PlayerID{
				"2",
				"3",
			},
			setup: func(myPoll *poll) {
				myPoll.ElectorIDs = []enum.PlayerID{
					"1",
					"2",
					"3",
				}
				myPoll.RemainingElectorIDs = []enum.PlayerID{
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
			assert.Equal(t, test.newElectorIDs, myPoll.(*poll).ElectorIDs)
			assert.Equal(t, test.newRemainingElectorIDs, myPoll.(*poll).RemainingElectorIDs)
		})
	}
}

func TestSetWeightPoll(t *testing.T) {
	tests := []struct {
		name           string
		electorID      enum.PlayerID
		weight         uint
		expectedStatus bool
		setup          func(*poll)
	}{
		{
			name:           "Failure (Non-existent elector)",
			electorID:      "99",
			weight:         1,
			expectedStatus: false,
			setup: func(myPoll *poll) {
				myPoll.RemainingElectorIDs = []enum.PlayerID{
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
				myPoll.RemainingElectorIDs = []enum.PlayerID{
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
		electorID      enum.PlayerID
		candidateID    enum.PlayerID
		expectedStatus bool
		newWeight      uint
		newVotes       uint
		setup          func(*poll)
	}{
		{
			name:           "Failure (Cannot vote)",
			electorID:      "99",
			candidateID:    "2",
			expectedStatus: false,
			setup: func(myPoll *poll) {
				// Open poll
				myPoll.Round = 1
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed: false,
				}

				myPoll.RemainingElectorIDs = []enum.PlayerID{"1"}
				myPoll.RemainingCandidateIDs = []enum.PlayerID{"2"}
			},
		},
		{
			name:           "Failure (Non-existent candidate)",
			electorID:      "1",
			candidateID:    "99",
			expectedStatus: false,
			setup: func(myPoll *poll) {
				// Open poll
				myPoll.Round = 1
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed: false,
				}

				myPoll.RemainingElectorIDs = []enum.PlayerID{"1"}
				myPoll.RemainingCandidateIDs = []enum.PlayerID{"2"}
			},
		},
		{
			name:           "Ok (Skip)",
			electorID:      "1",
			candidateID:    "",
			expectedStatus: true,
			newVotes:       1,
			newWeight:      1,
			setup: func(myPoll *poll) {
				// Open poll
				myPoll.Round = 1
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed:    false,
					VoteRecords: make(map[enum.PlayerID]*types.VoteRecord),
				}

				myPoll.RemainingElectorIDs = []enum.PlayerID{"1"}
				myPoll.RemainingCandidateIDs = []enum.PlayerID{"2"}
				myPoll.SetWeight("1", 8)
			},
		},
		{
			name:           "Ok (Voted)",
			electorID:      "1",
			candidateID:    "2",
			expectedStatus: true,
			newVotes:       1,
			newWeight:      8,
			setup: func(myPoll *poll) {
				// Open poll
				myPoll.Round = 1
				myPoll.Records[myPoll.Round] = &types.PollRecord{
					IsClosed:    false,
					VoteRecords: make(map[enum.PlayerID]*types.VoteRecord),
				}

				myPoll.RemainingElectorIDs = []enum.PlayerID{"1"}
				myPoll.RemainingCandidateIDs = []enum.PlayerID{"2"}
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
					test.newVotes,
					myPoll.Record(enum.LastRound).VoteRecords[test.candidateID].Votes,
				)
				assert.Equal(t,
					test.newWeight,
					myPoll.Record(enum.LastRound).VoteRecords[test.candidateID].Weights,
				)
				assert.Contains(
					t,
					myPoll.Record(enum.LastRound).VoteRecords[test.candidateID].ElectorIDs,
					test.electorID,
				)
				assert.Contains(t, myPoll.(*poll).VotedElectorIDs, test.electorID)
			}
		})
	}
}
