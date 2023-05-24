package logic

import (
	"errors"
	"testing"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/types"

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
	ps.Empty(p.(*poll).Weights)
	ps.NotNil(p.(*poll).Records)
	ps.Empty(p.(*poll).Records)
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
				p.Round = constants.FirstRound
				p.Records[p.Round] = &types.PollRecord{
					IsClosed: true,
				}
			},
		},
		{
			name:           "Is open",
			expectedStatus: true,
			setup: func(p *poll) {
				p.Round = constants.FirstRound
				p.Records[p.Round] = &types.PollRecord{
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
	electorId := types.PlayerId("1")
	tests := []struct {
		name           string
		electorId      types.PlayerId
		expectedStatus bool
		expectedErr    error
		setup          func(*poll)
	}{
		{
			name:           "Cannot vote (Poll was closed)",
			electorId:      electorId,
			expectedStatus: false,
			expectedErr:    errors.New("Wait for the next poll (0) ᕙ(⇀‸↼‶)ᕗ"),
			setup:          func(p *poll) {},
		},
		{
			name:           "Cannot vote (Not an elector)",
			electorId:      electorId,
			expectedStatus: false,
			expectedErr:    errors.New("You're not allowed to vote ノ(ジ)ー'"),
			setup: func(p *poll) {
				p.Round = constants.FirstRound
				p.Records[p.Round] = &types.PollRecord{
					IsClosed: false,
				}
			},
		},
		{
			name:           "Ok",
			electorId:      electorId,
			expectedStatus: true,
			setup: func(p *poll) {
				// Open poll
				p.Round = constants.FirstRound
				p.Records[p.Round] = &types.PollRecord{
					IsClosed: false,
				}

				p.RemainingElectorIds = append(
					p.RemainingElectorIds,
					electorId,
				)
			},
		},
		{
			name:           "Cannot vote (Already voted)",
			electorId:      electorId,
			expectedStatus: false,
			expectedErr:    errors.New("Wait for the next round ಠ_ಠ"),
			setup: func(p *poll) {
				p.Round = constants.FirstRound
				p.Records[p.Round] = &types.PollRecord{
					IsClosed: false,
				}

				p.RemainingElectorIds = append(
					p.RemainingElectorIds,
					electorId,
				)
				p.VotedElectorIds = append(
					p.VotedElectorIds,
					electorId,
				)
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			p := NewPoll()
			test.setup(p.(*poll))

			res, err := p.CanVote(test.electorId)

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
		givenRound     types.Round
		expectedRecord *types.PollRecord
		setup          func(*poll)
	}{
		{
			name:           "Nil (Non-existent round)",
			givenRound:     constants.FirstRound,
			expectedRecord: nil,
			setup:          func(p *poll) {},
		},
		{
			name:       "Ok",
			givenRound: constants.FirstRound,
			expectedRecord: &types.PollRecord{
				WinnerId: types.PlayerId("98"),
			},
			setup: func(p *poll) {
				p.Round = constants.FirstRound
				p.Records[p.Round] = &types.PollRecord{
					WinnerId: types.PlayerId("98"),
				}
			},
		},
		{
			name:       "Ok (Latest round)",
			givenRound: constants.ZeroRound,
			expectedRecord: &types.PollRecord{
				WinnerId: types.PlayerId("98"),
			},
			setup: func(p *poll) {
				p.Round = constants.SecondRound
				p.Records[p.Round] = &types.PollRecord{
					WinnerId: types.PlayerId("98"),
				}
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			p := NewPoll()
			test.setup(p.(*poll))

			record := p.Record(test.givenRound)

			ps.Equal(test.expectedRecord, record)
		})
	}
}

func (ps PollSuite) TestOpen() {
	tests := []struct {
		name           string
		expectedStatus bool
		expectedErr    error
		expectedRound  types.Round
		setup          func(*poll)
	}{
		{
			name:           "Failure (Already open)",
			expectedStatus: false,
			expectedErr:    errors.New("Poll is already open!"),
			expectedRound:  constants.FirstRound,
			setup: func(p *poll) {
				p.Round = constants.FirstRound
				p.Records[p.Round] = &types.PollRecord{
					IsClosed: false,
				}
			},
		},
		{
			name:           "Failure (Not enough candidates)",
			expectedStatus: false,
			expectedErr:    errors.New("Number of candidates is too small!"),
			expectedRound:  constants.FirstRound,
			setup: func(p *poll) {
				p.Round = constants.FirstRound
				p.Records[p.Round] = &types.PollRecord{
					IsClosed: true,
				}
				p.RemainingCandidateIds = []types.PlayerId{"1"}
			},
		},
		{
			name:           "Ok",
			expectedStatus: true,
			expectedRound:  constants.SecondRound,
			setup: func(p *poll) {
				p.Round = constants.FirstRound
				p.Records[p.Round] = &types.PollRecord{
					IsClosed: true,
				}
				p.RemainingCandidateIds = []types.PlayerId{"1", "2"}
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			p := NewPoll()
			test.setup(p.(*poll))

			ok, err := p.Open()

			ps.Equal(test.expectedStatus, ok)
			ps.Equal(test.expectedRound, p.(*poll).Round)
			ps.Equal(test.expectedErr, err)

			if test.expectedStatus == true {
				ps.False(p.Record(p.(*poll).Round).IsClosed)
			}
		})
	}
}

func (ps PollSuite) TestCurrentWinnerId() {
	tests := []struct {
		name            string
		expectedWinerId types.PlayerId
		setup           func(*poll)
	}{
		{
			name:            "No major vote",
			expectedWinerId: types.PlayerId(""),
			setup: func(p *poll) {
				p.Round = constants.FirstRound
				p.RemainingElectorIds = []types.PlayerId{
					"1", "2", "3", "4", "5",
				}
				p.Records[p.Round] = &types.PollRecord{
					VoteRecords: map[types.PlayerId]*types.VoteRecord{
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
			expectedWinerId: types.PlayerId(""),
			setup: func(p *poll) {
				p.Round = constants.FirstRound
				p.RemainingElectorIds = []types.PlayerId{
					"1", "2", "3", "4",
				}
				p.Records[p.Round] = &types.PollRecord{
					VoteRecords: map[types.PlayerId]*types.VoteRecord{
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
			expectedWinerId: types.PlayerId("1"),
			setup: func(p *poll) {
				p.Round = constants.FirstRound
				p.RemainingElectorIds = []types.PlayerId{
					"1", "2", "3", "4",
				}
				p.Records[p.Round] = &types.PollRecord{
					VoteRecords: map[types.PlayerId]*types.VoteRecord{
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

			winnerId := p.(*poll).currentWinnerId()

			ps.Equal(test.expectedWinerId, winnerId)
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
				p.Round = constants.FirstRound
				p.Records[p.Round] = &types.PollRecord{
					IsClosed: true,
				}
			},
		},
		{
			name:           "Ok",
			expectedStatus: true,
			expectedRecord: &types.PollRecord{
				IsClosed: true,
				WinnerId: types.PlayerId("1"),
				VoteRecords: map[types.PlayerId]*types.VoteRecord{
					"1": {
						ElectorIds: []types.PlayerId{"2", "3", "4"},
						Weights:    3,
						Votes:      3,
					},
					"2": {
						ElectorIds: []types.PlayerId{"1"},
						Weights:    1,
						Votes:      1,
					},
					"": {
						ElectorIds: []types.PlayerId{"5"},
						Votes:      1,
						Weights:    1,
					},
				},
			},
			setup: func(p *poll) {
				p.Round = constants.FirstRound
				p.RemainingElectorIds = []types.PlayerId{
					"1",
					"2",
					"3",
					"4",
					"5",
				}
				p.VotedElectorIds = []types.PlayerId{
					"1",
					"2",
					"3",
					"4",
				}
				p.Records[p.Round] = &types.PollRecord{
					IsClosed: false,
					VoteRecords: map[types.PlayerId]*types.VoteRecord{
						"1": {
							ElectorIds: []types.PlayerId{"2", "3", "4"},
							Weights:    3,
							Votes:      3,
						},
						"2": {
							ElectorIds: []types.PlayerId{"1"},
							Weights:    1,
							Votes:      1,
						},
					},
				}
				p.Records[p.Round].VoteRecords[""] = &types.VoteRecord{
					ElectorIds: []types.PlayerId{},
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
			ps.Equal(test.expectedRecord, p.(*poll).Records[p.(*poll).Round])

			if test.expectedStatus == true {
				ps.False(p.IsOpen())
				ps.Empty(p.(*poll).VotedElectorIds)
			}
		})
	}
}

func (ps PollSuite) TestAddCandidates() {
	tests := []struct {
		name                     string
		candIdateIds             []types.PlayerId
		newCandidateIds          []types.PlayerId
		newRemainingCandIdateIds []types.PlayerId
		setup                    func(*poll)
	}{
		{
			name: "Failure (Already existed in remaining)",
			candIdateIds: []types.PlayerId{
				"1",
				"2",
			},
			newCandidateIds: []types.PlayerId{
				"1",
				"2",
				"3",
			},
			newRemainingCandIdateIds: []types.PlayerId{
				"1",
				"2",
				"3",
			},
			setup: func(p *poll) {
				p.CandidateIds = []types.PlayerId{
					"1",
					"2",
					"3",
				}
				p.RemainingCandidateIds = []types.PlayerId{
					"1",
					"2",
					"3",
				}
			},
		},
		{
			name: "Ok (Doesn't exist in remaining but in all)",
			candIdateIds: []types.PlayerId{
				"1",
				"2",
			},
			newCandidateIds: []types.PlayerId{
				"1",
				"2",
				"3",
			},
			newRemainingCandIdateIds: []types.PlayerId{
				"2",
				"3",
				"1",
			},
			setup: func(p *poll) {
				p.CandidateIds = []types.PlayerId{
					"1",
					"2",
					"3",
				}
				p.RemainingCandidateIds = []types.PlayerId{
					"2",
					"3",
				}
			},
		},
		{
			name: "Ok (Doesn't exist in remaining and all)",
			candIdateIds: []types.PlayerId{
				"1",
				"2",
			},
			newCandidateIds: []types.PlayerId{
				"2",
				"3",
				"1",
			},
			newRemainingCandIdateIds: []types.PlayerId{
				"2",
				"3",
				"1",
			},
			setup: func(p *poll) {
				p.CandidateIds = []types.PlayerId{
					"2",
					"3",
				}
				p.RemainingCandidateIds = []types.PlayerId{
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

			p.AddCandidates(test.candIdateIds...)

			ps.Equal(test.newCandidateIds, p.(*poll).CandidateIds)
			ps.Equal(test.newRemainingCandIdateIds, p.(*poll).RemainingCandidateIds)
		})
	}
}

func (ps PollSuite) TestRemoveCandidate() {
	tests := []struct {
		name                     string
		candIdateId              types.PlayerId
		expectedStatus           bool
		newCandidateIds          []types.PlayerId
		newRemainingCandIdateIds []types.PlayerId
		setup                    func(*poll)
	}{
		{
			name:           "Failure (Non-existent candIdate)",
			candIdateId:    types.PlayerId("99"),
			expectedStatus: false,
			newCandidateIds: []types.PlayerId{
				"1",
				"2",
				"3",
			},
			newRemainingCandIdateIds: []types.PlayerId{
				"1",
				"2",
				"3",
			},
			setup: func(p *poll) {
				p.CandidateIds = []types.PlayerId{
					"1",
					"2",
					"3",
				}
				p.RemainingCandidateIds = []types.PlayerId{
					"1",
					"2",
					"3",
				}
			},
		},
		{
			name:           "Ok",
			candIdateId:    types.PlayerId("1"),
			expectedStatus: true,
			newCandidateIds: []types.PlayerId{
				"1",
				"2",
				"3",
			},
			newRemainingCandIdateIds: []types.PlayerId{
				"2",
				"3",
			},
			setup: func(p *poll) {
				p.CandidateIds = []types.PlayerId{
					"1",
					"2",
					"3",
				}
				p.RemainingCandidateIds = []types.PlayerId{
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

			ok := p.RemoveCandidate(test.candIdateId)

			ps.Equal(test.expectedStatus, ok)
			ps.Equal(test.newCandidateIds, p.(*poll).CandidateIds)
			ps.Equal(test.newRemainingCandIdateIds, p.(*poll).RemainingCandidateIds)
		})
	}
}

func (ps PollSuite) TestAddElectors() {
	tests := []struct {
		name                   string
		electorIds             []types.PlayerId
		newElectorIds          []types.PlayerId
		newRemainingElectorIds []types.PlayerId
		setup                  func(*poll)
	}{
		{
			name: "Already existed in remaining",
			electorIds: []types.PlayerId{
				"1",
				"2",
			},
			newElectorIds: []types.PlayerId{
				"1",
				"2",
				"3",
			},
			newRemainingElectorIds: []types.PlayerId{
				"1",
				"2",
				"3",
			},
			setup: func(p *poll) {
				p.ElectorIds = []types.PlayerId{
					"1",
					"2",
					"3",
				}
				p.RemainingElectorIds = []types.PlayerId{
					"1",
					"2",
					"3",
				}
			},
		},
		{
			name: "Doesn't exist in remaining but in all",
			electorIds: []types.PlayerId{
				"1",
				"2",
			},
			newElectorIds: []types.PlayerId{
				"1",
				"2",
				"3",
			},
			newRemainingElectorIds: []types.PlayerId{
				"2",
				"3",
				"1",
			},
			setup: func(p *poll) {
				p.ElectorIds = []types.PlayerId{
					"1",
					"2",
					"3",
				}
				p.RemainingElectorIds = []types.PlayerId{
					"2",
					"3",
				}
			},
		},
		{
			name: "Doesn't exist in remaining and all",
			electorIds: []types.PlayerId{
				"1",
				"2",
			},
			newElectorIds: []types.PlayerId{
				"2",
				"3",
				"1",
			},
			newRemainingElectorIds: []types.PlayerId{
				"2",
				"3",
				"1",
			},
			setup: func(p *poll) {
				p.ElectorIds = []types.PlayerId{
					"2",
					"3",
				}
				p.RemainingElectorIds = []types.PlayerId{
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

			p.AddElectors(test.electorIds...)

			ps.Equal(test.newElectorIds, p.(*poll).ElectorIds)
			ps.Equal(test.newRemainingElectorIds, p.(*poll).RemainingElectorIds)
		})
	}
}

func (ps PollSuite) TestRemoveElector() {
	tests := []struct {
		name                   string
		electorId              types.PlayerId
		expectedStatus         bool
		newElectorIds          []types.PlayerId
		newRemainingElectorIds []types.PlayerId
		setup                  func(*poll)
	}{
		{
			name:           "Failure (Non-existent elector)",
			electorId:      types.PlayerId("99"),
			expectedStatus: false,
			newElectorIds: []types.PlayerId{
				"1",
				"2",
				"3",
			},
			newRemainingElectorIds: []types.PlayerId{
				"1",
				"2",
				"3",
			},
			setup: func(p *poll) {
				p.ElectorIds = []types.PlayerId{
					"1",
					"2",
					"3",
				}
				p.RemainingElectorIds = []types.PlayerId{
					"1",
					"2",
					"3",
				}
			},
		},
		{
			name:           "Ok",
			electorId:      types.PlayerId("1"),
			expectedStatus: true,
			newElectorIds: []types.PlayerId{
				"1",
				"2",
				"3",
			},
			newRemainingElectorIds: []types.PlayerId{
				"2",
				"3",
			},
			setup: func(p *poll) {
				p.ElectorIds = []types.PlayerId{
					"1",
					"2",
					"3",
				}
				p.RemainingElectorIds = []types.PlayerId{
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

			ok := p.RemoveElector(test.electorId)

			ps.Equal(test.expectedStatus, ok)
			ps.Equal(test.newElectorIds, p.(*poll).ElectorIds)
			ps.Equal(test.newRemainingElectorIds, p.(*poll).RemainingElectorIds)
		})
	}
}

func (ps PollSuite) TestSetWeight() {
	tests := []struct {
		name           string
		electorId      types.PlayerId
		weight         uint
		expectedStatus bool
		setup          func(*poll)
	}{
		{
			name:           "Failure (Non-existent elector)",
			electorId:      types.PlayerId("99"),
			weight:         1,
			expectedStatus: false,
			setup: func(p *poll) {
				p.RemainingElectorIds = []types.PlayerId{
					"1",
					"2",
				}
			},
		},
		{
			name:           "Ok",
			electorId:      types.PlayerId("1"),
			weight:         5,
			expectedStatus: true,
			setup: func(p *poll) {
				p.RemainingElectorIds = []types.PlayerId{
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

			ok := p.SetWeight(test.electorId, test.weight)

			ps.Equal(test.expectedStatus, ok)
			if test.expectedStatus == true {
				ps.Equal(p.(*poll).Weights[test.electorId], test.weight)
			}
		})
	}
}

func (ps PollSuite) TestVote() {
	tests := []struct {
		name           string
		electorId      types.PlayerId
		candIdateId    types.PlayerId
		expectedStatus bool
		expectedErr    error
		newWeight      uint
		newVotes       uint
		setup          func(*poll)
	}{
		{
			name:           "Failure (Cannot vote)",
			electorId:      types.PlayerId("99"),
			candIdateId:    types.PlayerId("2"),
			expectedStatus: false,
			expectedErr:    errors.New("You're not allowed to vote ノ(ジ)ー'"),
			setup: func(p *poll) {
				p.Round = constants.FirstRound
				p.Records[p.Round] = &types.PollRecord{
					IsClosed: false,
				}

				p.RemainingElectorIds = []types.PlayerId{"1"}
				p.RemainingCandidateIds = []types.PlayerId{"2"}
			},
		},
		{
			name:           "Failure (Non-existent candIdate)",
			electorId:      types.PlayerId("1"),
			candIdateId:    types.PlayerId("99"),
			expectedStatus: false,
			expectedErr:    errors.New("Your vote is not valid ¬_¬"),
			setup: func(p *poll) {
				p.Round = constants.FirstRound
				p.Records[p.Round] = &types.PollRecord{
					IsClosed: false,
				}

				p.RemainingElectorIds = []types.PlayerId{"1"}
				p.RemainingCandidateIds = []types.PlayerId{"2"}
			},
		},
		{
			name:           "Ok (Skip)",
			electorId:      types.PlayerId("1"),
			candIdateId:    types.PlayerId(""),
			expectedStatus: true,
			newVotes:       1,
			newWeight:      1,
			setup: func(p *poll) {
				p.Round = constants.FirstRound
				p.Records[p.Round] = &types.PollRecord{
					IsClosed:    false,
					VoteRecords: make(map[types.PlayerId]*types.VoteRecord),
				}

				p.RemainingElectorIds = []types.PlayerId{"1"}
				p.RemainingCandidateIds = []types.PlayerId{"2"}
				p.SetWeight("1", 8)
			},
		},
		{
			name:           "Ok (Voted)",
			electorId:      types.PlayerId("1"),
			candIdateId:    types.PlayerId("2"),
			expectedStatus: true,
			newVotes:       1,
			newWeight:      8,
			setup: func(p *poll) {
				p.Round = constants.FirstRound
				p.Records[p.Round] = &types.PollRecord{
					IsClosed:    false,
					VoteRecords: make(map[types.PlayerId]*types.VoteRecord),
				}

				p.RemainingElectorIds = []types.PlayerId{"1"}
				p.RemainingCandidateIds = []types.PlayerId{"2"}
				p.SetWeight("1", 8)
			},
		},
	}

	for _, test := range tests {
		ps.Run(test.name, func() {
			p := NewPoll()
			test.setup(p.(*poll))

			ok, err := p.Vote(test.electorId, test.candIdateId)

			ps.Equal(test.expectedStatus, ok)
			if test.expectedStatus == true {
				ps.Equal(
					test.newVotes,
					p.Record(p.(*poll).Round).VoteRecords[test.candIdateId].Votes,
				)
				ps.Equal(
					test.newWeight,
					p.Record(p.(*poll).Round).VoteRecords[test.candIdateId].Weights,
				)
				ps.Contains(
					p.Record(p.(*poll).Round).VoteRecords[test.candIdateId].ElectorIds,
					test.electorId,
				)
				ps.Contains(p.(*poll).VotedElectorIds, test.electorId)
			} else {
				ps.Equal(test.expectedErr, err)
			}
		})
	}
}
