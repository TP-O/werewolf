package core

import (
	"fmt"
	"testing"
	"uwwolf/game/types"
	"uwwolf/game/vars"

	"github.com/stretchr/testify/suite"
)

type PollSuite struct {
	suite.Suite
}

func TestPollSuite(t *testing.T) {
	suite.Run(t, new(PollSuite))
}

func (ps PollSuite) TestNewPoll() {
	p := NewPoll()

	ps.NotNil(p)
	ps.NotNil(p.(*poll).Weights)
	ps.Len(p.(*poll).Weights, 0)
	ps.NotNil(p.(*poll).Records)
	ps.Len(p.(*poll).Records, 0)
}

func (ps PollSuite) TestIsOpen() {
	tests := []struct {
		name           string
		expectedStatus bool
		setup          func(*poll)
	}{
		{
			name:           "Isnt open (zero round)",
			expectedStatus: false,
			setup:          func(p *poll) {},
		},
		{
			name:           "Isnt open (closed)",
			expectedStatus: false,
			setup: func(p *poll) {
				p.RoundID = vars.FirstRound
				p.Records[p.RoundID] = &types.PollRecord{
					IsClosed: true,
				}
			},
		},
		{
			name:           "Is open",
			expectedStatus: true,
			setup: func(p *poll) {
				p.RoundID = vars.FirstRound
				p.Records[p.RoundID] = &types.PollRecord{
					IsClosed: false,
				}
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			p := NewPoll()
			test.setup(p.(*poll))

			isOpen := p.IsOpen()

			ps.Equal(test.expectedStatus, isOpen)
		})
	}
}

func (ps PollSuite) TestCanVote() {
	electorID := types.PlayerID("1")
	tests := []struct {
		name           string
		electorID      types.PlayerID
		expectedStatus bool
		expectedErr    error
		setup          func(*poll)
	}{
		{
			name:           "Cannot vote (Poll was closed)",
			electorID:      electorID,
			expectedStatus: false,
			expectedErr:    fmt.Errorf("Wait for the next poll (%v) ᕙ(⇀‸↼‶)ᕗ", 0),
			setup:          func(p *poll) {},
		},
		{
			name:           "Cannot vote (Not an elector)",
			electorID:      electorID,
			expectedStatus: false,
			expectedErr:    fmt.Errorf("You're not allowed to vote ノ(ジ)ー'"),
			setup: func(p *poll) {
				p.RoundID = vars.FirstRound
				p.Records[p.RoundID] = &types.PollRecord{
					IsClosed: false,
				}
			},
		},
		{
			name:           "Ok",
			electorID:      electorID,
			expectedStatus: true,
			setup: func(p *poll) {
				// Open poll
				p.RoundID = vars.FirstRound
				p.Records[p.RoundID] = &types.PollRecord{
					IsClosed: false,
				}

				p.RemainingElectorIDs = append(
					p.RemainingElectorIDs,
					electorID,
				)
			},
		},
		{
			name:           "Cannot vote (Already voted)",
			electorID:      electorID,
			expectedStatus: false,
			expectedErr:    fmt.Errorf("Wait for the next round ಠ_ಠ"),
			setup: func(p *poll) {
				p.RoundID = vars.FirstRound
				p.Records[p.RoundID] = &types.PollRecord{
					IsClosed: false,
				}

				p.RemainingElectorIDs = append(
					p.RemainingElectorIDs,
					electorID,
				)
				p.VotedElectorIDs = append(
					p.VotedElectorIDs,
					electorID,
				)
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			p := NewPoll()
			test.setup(p.(*poll))

			res, err := p.CanVote(test.electorID)

			ps.Equal(test.expectedStatus, res)
			if res == false {
				ps.Equal(test.expectedErr, err)
			}
		})
	}
}

func (ps PollSuite) TestRecord() {
	tests := []struct {
		name           string
		givenRoundID   types.RoundID
		expectedRecord *types.PollRecord
		setup          func(*poll)
	}{
		{
			name:           "Nil (Non-existent round)",
			givenRoundID:   vars.FirstRound,
			expectedRecord: nil,
			setup:          func(p *poll) {},
		},
		{
			name:         "Ok",
			givenRoundID: vars.FirstRound,
			expectedRecord: &types.PollRecord{
				WinnerID: types.PlayerID("98"),
			},
			setup: func(p *poll) {
				p.RoundID = vars.FirstRound
				p.Records[p.RoundID] = &types.PollRecord{
					WinnerID: types.PlayerID("98"),
				}
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			p := NewPoll()
			test.setup(p.(*poll))

			record := p.Record(test.givenRoundID)

			ps.Equal(test.expectedRecord, record)
		})
	}
}

func (ps PollSuite) TestOpen() {
	tests := []struct {
		name            string
		expectedStatus  bool
		expectedErr     error
		expectedRoundID types.RoundID
		setup           func(*poll)
	}{
		{
			name:            "Failure (Already open)",
			expectedStatus:  false,
			expectedErr:     fmt.Errorf("Poll is already open!"),
			expectedRoundID: vars.FirstRound,
			setup: func(p *poll) {
				p.RoundID = vars.FirstRound
				p.Records[p.RoundID] = &types.PollRecord{
					IsClosed: false,
				}
			},
		},
		{
			name:            "Failure (Not enough candidates)",
			expectedStatus:  false,
			expectedErr:     fmt.Errorf("Number of candidates is too small!"),
			expectedRoundID: vars.FirstRound,
			setup: func(p *poll) {
				p.RoundID = vars.FirstRound
				p.Records[p.RoundID] = &types.PollRecord{
					IsClosed: true,
				}
				p.RemainingCandidateIDs = []types.PlayerID{"1"}
			},
		},
		{
			name:            "Ok",
			expectedStatus:  true,
			expectedRoundID: vars.SecondRound,
			setup: func(p *poll) {
				p.RoundID = vars.FirstRound
				p.Records[p.RoundID] = &types.PollRecord{
					IsClosed: true,
				}
				p.RemainingCandidateIDs = []types.PlayerID{"1", "2"}
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			p := NewPoll()
			test.setup(p.(*poll))

			ok, err := p.Open()

			ps.Equal(test.expectedStatus, ok)
			ps.Equal(test.expectedRoundID, p.(*poll).RoundID)
			ps.Equal(test.expectedErr, err)

			if test.expectedStatus == true {
				ps.False(p.Record(p.(*poll).RoundID).IsClosed)
			}
		})
	}
}

func (ps PollSuite) TestCurrentWinnerID() {
	tests := []struct {
		name            string
		expectedWinerID types.PlayerID
		setup           func(*poll)
	}{
		{
			name:            "No major vote",
			expectedWinerID: types.PlayerID(""),
			setup: func(p *poll) {
				p.RoundID = vars.FirstRound
				p.RemainingElectorIDs = []types.PlayerID{
					"1", "2", "3", "4", "5",
				}
				p.Records[p.RoundID] = &types.PollRecord{
					VoteRecords: map[types.PlayerID]*types.VoteRecord{
						"1": {
							Weights: 2,
						},
						"2": {
							Weights: 2,
						},
						"": {
							Weights: 1,
						},
					},
				}
			},
		},
		{
			name:            "Draw vote",
			expectedWinerID: types.PlayerID(""),
			setup: func(p *poll) {
				p.RoundID = vars.FirstRound
				p.RemainingElectorIDs = []types.PlayerID{
					"1", "2", "3", "4",
				}
				p.Records[p.RoundID] = &types.PollRecord{
					VoteRecords: map[types.PlayerID]*types.VoteRecord{
						"1": {
							Weights: 2,
						},
						"2": {
							Weights: 2,
						},
					},
				}
			},
		},
		{
			name:            "Major vote",
			expectedWinerID: types.PlayerID("1"),
			setup: func(p *poll) {
				p.RoundID = vars.FirstRound
				p.RemainingElectorIDs = []types.PlayerID{
					"1", "2", "3", "4",
				}
				p.Records[p.RoundID] = &types.PollRecord{
					VoteRecords: map[types.PlayerID]*types.VoteRecord{
						"1": {
							Weights: 3,
						},
						"2": {
							Weights: 1,
						},
					},
				}
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			p := NewPoll()
			test.setup(p.(*poll))

			winnerID := p.(*poll).currentWinnerID()

			ps.Equal(test.expectedWinerID, winnerID)
		})
	}
}

func (ps PollSuite) TestClose() {
	tests := []struct {
		name           string
		expectedStatus bool
		expectedRecord *types.PollRecord
		setup          func(*poll)
	}{
		{
			name:           "Failure (Poll was closed)",
			expectedStatus: false,
			expectedRecord: &types.PollRecord{
				IsClosed: true,
			},
			setup: func(p *poll) {
				p.RoundID = vars.FirstRound
				p.Records[p.RoundID] = &types.PollRecord{
					IsClosed: true,
				}
			},
		},
		{
			name:           "Ok",
			expectedStatus: true,
			expectedRecord: &types.PollRecord{
				IsClosed: true,
				WinnerID: types.PlayerID("1"),
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
			setup: func(p *poll) {
				p.RoundID = vars.FirstRound
				p.RemainingElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
					"4",
					"5",
				}
				p.VotedElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
					"4",
				}
				p.Records[p.RoundID] = &types.PollRecord{
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
				p.Records[p.RoundID].VoteRecords[""] = &types.VoteRecord{
					ElectorIDs: []types.PlayerID{},
				}
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			p := NewPoll()
			test.setup(p.(*poll))

			ok := p.Close()

			ps.Equal(test.expectedStatus, ok)
			ps.Equal(test.expectedRecord, p.(*poll).Records[p.(*poll).RoundID])

			if test.expectedStatus == true {
				ps.False(p.IsOpen())
				ps.Empty(p.(*poll).VotedElectorIDs)
			}
		})
	}
}

func (ps PollSuite) TestAddCandidates() {
	tests := []struct {
		name                     string
		candidateIDs             []types.PlayerID
		newCandidateIDs          []types.PlayerID
		newRemainingCandidateIDs []types.PlayerID
		setup                    func(*poll)
	}{
		{
			name: "Failure (Already existed in remaining)",
			candidateIDs: []types.PlayerID{
				"1",
				"2",
			},
			newCandidateIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			newRemainingCandidateIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			setup: func(p *poll) {
				p.CandidateIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
				p.RemainingCandidateIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
			},
		},
		{
			name: "Ok (Doesn't exist in remaining but in all)",
			candidateIDs: []types.PlayerID{
				"1",
				"2",
			},
			newCandidateIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			newRemainingCandidateIDs: []types.PlayerID{
				"2",
				"3",
				"1",
			},
			setup: func(p *poll) {
				p.CandidateIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
				p.RemainingCandidateIDs = []types.PlayerID{
					"2",
					"3",
				}
			},
		},
		{
			name: "Ok (Doesn't exist in remaining and all)",
			candidateIDs: []types.PlayerID{
				"1",
				"2",
			},
			newCandidateIDs: []types.PlayerID{
				"2",
				"3",
				"1",
			},
			newRemainingCandidateIDs: []types.PlayerID{
				"2",
				"3",
				"1",
			},
			setup: func(p *poll) {
				p.CandidateIDs = []types.PlayerID{
					"2",
					"3",
				}
				p.RemainingCandidateIDs = []types.PlayerID{
					"2",
					"3",
				}
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			p := NewPoll()
			test.setup(p.(*poll))

			p.AddCandidates(test.candidateIDs...)

			ps.Equal(test.newCandidateIDs, p.(*poll).CandidateIDs)
			ps.Equal(test.newRemainingCandidateIDs, p.(*poll).RemainingCandidateIDs)
		})
	}
}

func (ps PollSuite) TestRemoveCandidate() {
	tests := []struct {
		name                     string
		candidateID              types.PlayerID
		expectedStatus           bool
		newCandidateIDs          []types.PlayerID
		newRemainingCandidateIDs []types.PlayerID
		setup                    func(*poll)
	}{
		{
			name:           "Failure (Non-existent candidate)",
			candidateID:    types.PlayerID("99"),
			expectedStatus: false,
			newCandidateIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			newRemainingCandidateIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			setup: func(p *poll) {
				p.CandidateIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
				p.RemainingCandidateIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
			},
		},
		{
			name:           "Ok",
			candidateID:    types.PlayerID("1"),
			expectedStatus: true,
			newCandidateIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			newRemainingCandidateIDs: []types.PlayerID{
				"2",
				"3",
			},
			setup: func(p *poll) {
				p.CandidateIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
				p.RemainingCandidateIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			p := NewPoll()
			test.setup(p.(*poll))

			ok := p.RemoveCandidate(test.candidateID)

			ps.Equal(test.expectedStatus, ok)
			ps.Equal(test.newCandidateIDs, p.(*poll).CandidateIDs)
			ps.Equal(test.newRemainingCandidateIDs, p.(*poll).RemainingCandidateIDs)
		})
	}
}

func (ps PollSuite) TestAddElectors() {
	tests := []struct {
		name                   string
		electorIDs             []types.PlayerID
		newElectorIDs          []types.PlayerID
		newRemainingElectorIDs []types.PlayerID
		setup                  func(*poll)
	}{
		{
			name: "Already existed in remaining",
			electorIDs: []types.PlayerID{
				"1",
				"2",
			},
			newElectorIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			newRemainingElectorIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			setup: func(p *poll) {
				p.ElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
				p.RemainingElectorIDs = []types.PlayerID{
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
			newElectorIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			newRemainingElectorIDs: []types.PlayerID{
				"2",
				"3",
				"1",
			},
			setup: func(p *poll) {
				p.ElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
				p.RemainingElectorIDs = []types.PlayerID{
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
			newElectorIDs: []types.PlayerID{
				"2",
				"3",
				"1",
			},
			newRemainingElectorIDs: []types.PlayerID{
				"2",
				"3",
				"1",
			},
			setup: func(p *poll) {
				p.ElectorIDs = []types.PlayerID{
					"2",
					"3",
				}
				p.RemainingElectorIDs = []types.PlayerID{
					"2",
					"3",
				}
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			p := NewPoll()
			test.setup(p.(*poll))

			p.AddElectors(test.electorIDs...)

			ps.Equal(test.newElectorIDs, p.(*poll).ElectorIDs)
			ps.Equal(test.newRemainingElectorIDs, p.(*poll).RemainingElectorIDs)
		})
	}
}

func (ps PollSuite) TestRemoveElector() {
	tests := []struct {
		name                   string
		electorID              types.PlayerID
		expectedStatus         bool
		newElectorIDs          []types.PlayerID
		newRemainingElectorIDs []types.PlayerID
		setup                  func(*poll)
	}{
		{
			name:           "Failure (Non-existent elector)",
			electorID:      types.PlayerID("99"),
			expectedStatus: false,
			newElectorIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			newRemainingElectorIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			setup: func(p *poll) {
				p.ElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
				p.RemainingElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
			},
		},
		{
			name:           "Ok",
			electorID:      types.PlayerID("1"),
			expectedStatus: true,
			newElectorIDs: []types.PlayerID{
				"1",
				"2",
				"3",
			},
			newRemainingElectorIDs: []types.PlayerID{
				"2",
				"3",
			},
			setup: func(p *poll) {
				p.ElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
				p.RemainingElectorIDs = []types.PlayerID{
					"1",
					"2",
					"3",
				}
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			p := NewPoll()
			test.setup(p.(*poll))

			ok := p.RemoveElector(test.electorID)

			ps.Equal(test.expectedStatus, ok)
			ps.Equal(test.newElectorIDs, p.(*poll).ElectorIDs)
			ps.Equal(test.newRemainingElectorIDs, p.(*poll).RemainingElectorIDs)
		})
	}
}

func (ps PollSuite) TestSetWeight() {
	tests := []struct {
		name           string
		electorID      types.PlayerID
		weight         uint
		expectedStatus bool
		setup          func(*poll)
	}{
		{
			name:           "Failure (Non-existent elector)",
			electorID:      types.PlayerID("99"),
			weight:         1,
			expectedStatus: false,
			setup: func(p *poll) {
				p.RemainingElectorIDs = []types.PlayerID{
					"1",
					"2",
				}
			},
		},
		{
			name:           "Ok",
			electorID:      types.PlayerID("1"),
			weight:         5,
			expectedStatus: true,
			setup: func(p *poll) {
				p.RemainingElectorIDs = []types.PlayerID{
					"1",
					"2",
				}
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			p := NewPoll()
			test.setup(p.(*poll))

			ok := p.SetWeight(test.electorID, test.weight)

			ps.Equal(test.expectedStatus, ok)
			if test.expectedStatus == true {
				ps.Equal(p.(*poll).Weights[test.electorID], test.weight)
			}
		})
	}
}

func (ps PollSuite) TestVote() {
	tests := []struct {
		name           string
		electorID      types.PlayerID
		candidateID    types.PlayerID
		expectedStatus bool
		expectedErr    error
		newWeight      uint
		newVotes       uint
		setup          func(*poll)
	}{
		{
			name:           "Failure (Cannot vote)",
			electorID:      types.PlayerID("99"),
			candidateID:    types.PlayerID("2"),
			expectedStatus: false,
			expectedErr:    fmt.Errorf("You're not allowed to vote ノ(ジ)ー'"),
			setup: func(p *poll) {
				p.RoundID = vars.FirstRound
				p.Records[p.RoundID] = &types.PollRecord{
					IsClosed: false,
				}

				p.RemainingElectorIDs = []types.PlayerID{"1"}
				p.RemainingCandidateIDs = []types.PlayerID{"2"}
			},
		},
		{
			name:           "Failure (Non-existent candidate)",
			electorID:      types.PlayerID("1"),
			candidateID:    types.PlayerID("99"),
			expectedStatus: false,
			expectedErr:    fmt.Errorf("Your vote is not valid ¬_¬"),
			setup: func(p *poll) {
				p.RoundID = vars.FirstRound
				p.Records[p.RoundID] = &types.PollRecord{
					IsClosed: false,
				}

				p.RemainingElectorIDs = []types.PlayerID{"1"}
				p.RemainingCandidateIDs = []types.PlayerID{"2"}
			},
		},
		{
			name:           "Ok (Skip)",
			electorID:      types.PlayerID("1"),
			candidateID:    types.PlayerID(""),
			expectedStatus: true,
			newVotes:       1,
			newWeight:      1,
			setup: func(p *poll) {
				p.RoundID = vars.FirstRound
				p.Records[p.RoundID] = &types.PollRecord{
					IsClosed:    false,
					VoteRecords: make(map[types.PlayerID]*types.VoteRecord),
				}

				p.RemainingElectorIDs = []types.PlayerID{"1"}
				p.RemainingCandidateIDs = []types.PlayerID{"2"}
				p.SetWeight("1", 8)
			},
		},
		{
			name:           "Ok (Voted)",
			electorID:      types.PlayerID("1"),
			candidateID:    types.PlayerID("2"),
			expectedStatus: true,
			newVotes:       1,
			newWeight:      8,
			setup: func(p *poll) {
				p.RoundID = vars.FirstRound
				p.Records[p.RoundID] = &types.PollRecord{
					IsClosed:    false,
					VoteRecords: make(map[types.PlayerID]*types.VoteRecord),
				}

				p.RemainingElectorIDs = []types.PlayerID{"1"}
				p.RemainingCandidateIDs = []types.PlayerID{"2"}
				p.SetWeight("1", 8)
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			p := NewPoll()
			test.setup(p.(*poll))

			ok, err := p.Vote(test.electorID, test.candidateID)

			ps.Equal(test.expectedStatus, ok)
			if test.expectedStatus == true {
				ps.Equal(
					test.newVotes,
					p.Record(p.(*poll).RoundID).VoteRecords[test.candidateID].Votes,
				)
				ps.Equal(
					test.newWeight,
					p.Record(p.(*poll).RoundID).VoteRecords[test.candidateID].Weights,
				)
				ps.Contains(
					p.Record(p.(*poll).RoundID).VoteRecords[test.candidateID].ElectorIDs,
					test.electorID,
				)
				ps.Contains(p.(*poll).VotedElectorIDs, test.electorID)
			} else {
				ps.Equal(test.expectedErr, err)
			}
		})
	}
}
