package projection

import (
	"testing"
	"time"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/org"
)

func TestFeatureProjection_reduces(t *testing.T) {
	type args struct {
		event func(t *testing.T) eventstore.EventReader
	}
	tests := []struct {
		name   string
		args   args
		reduce func(event eventstore.EventReader) (*handler.Statement, error)
		want   wantReduce
	}{
		{
			name: "org.reduceFeatureSet",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.FeaturesSetEventType),
					org.AggregateType,
					[]byte(`{
						"tierName": "TierName",
						"tierDescription": "TierDescription",
						"state": 1,
						"stateDescription": "StateDescription",
						"auditLogRetention": 1,
						"loginPolicyFactors": true,
						"loginPolicyIDP": true,
						"loginPolicyPasswordless": true,
						"loginPolicyRegistration": true,
						"loginPolicyUsernameLogin": true,
						"loginPolicyPasswordReset": true,
						"passwordComplexityPolicy": true,
						"labelPolicyPrivateLabel": true,
						"labelPolicyWatermark": true,
						"customDomain": true,
						"privacyPolicy": true,
						"metadataUser": true,
						"customTextMessage": true,
						"customTextLogin": true,
						"lockoutPolicy": true,
						"actions": true
					}`),
				), org.FeaturesSetEventMapper),
			},
			reduce: (&FeatureProjection{}).reduceFeatureSet,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       FeatureTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPSERT INTO zitadel.projections.features (aggregate_id, creation_date, change_date, sequence, is_default, tier_name, tier_description, state, state_description, audit_log_retention, login_policy_factors, login_policy_idp, login_policy_passwordless, login_policy_registration, login_policy_username_login, login_policy_password_reset, password_complexity_policy, label_policy_private_label, label_policy_watermark, custom_domain, privacy_policy, metadata_user, custom_text_message, custom_text_login, lockout_policy, actions) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								uint64(15),
								false,
								"TierName",
								"TierDescription",
								domain.FeaturesStateActive,
								"StateDescription",
								time.Nanosecond,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceFeatureRemoved",
			reduce: (&FeatureProjection{}).reduceFeatureRemoved,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.FeaturesRemovedEventType),
					org.AggregateType,
					nil,
				), org.FeaturesRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       FeatureTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.features WHERE (aggregate_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduceFeatureSet",
			reduce: (&FeatureProjection{}).reduceFeatureSet,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.FeaturesSetEventType),
					iam.AggregateType,
					[]byte(`{
				"tierName": "TierName",
				"tierDescription": "TierDescription",
				"state": 1,
				"stateDescription": "StateDescription",
				"auditLogRetention": 1,
				"loginPolicyFactors": true,
				"loginPolicyIDP": true,
				"loginPolicyPasswordless": true,
				"loginPolicyRegistration": true,
				"loginPolicyUsernameLogin": true,
				"loginPolicyPasswordReset": true,
				"passwordComplexityPolicy": true,
				"labelPolicyPrivateLabel": true,
				"labelPolicyWatermark": true,
				"customDomain": true,
				"privacyPolicy": true,
				"metadataUser": true,
				"customTextMessage": true,
				"customTextLogin": true,
				"lockoutPolicy": true,
				"actions": true
			}`),
				), iam.FeaturesSetEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       FeatureTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPSERT INTO zitadel.projections.features (aggregate_id, creation_date, change_date, sequence, is_default, tier_name, tier_description, state, state_description, audit_log_retention, login_policy_factors, login_policy_idp, login_policy_passwordless, login_policy_registration, login_policy_username_login, login_policy_password_reset, password_complexity_policy, label_policy_private_label, label_policy_watermark, custom_domain, privacy_policy, metadata_user, custom_text_message, custom_text_login, lockout_policy, actions) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								uint64(15),
								true,
								"TierName",
								"TierDescription",
								domain.FeaturesStateActive,
								"StateDescription",
								time.Nanosecond,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
								true,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := baseEvent(t)
			got, err := tt.reduce(event)
			if _, ok := err.(errors.InvalidArgument); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, tt.want)
		})
	}
}