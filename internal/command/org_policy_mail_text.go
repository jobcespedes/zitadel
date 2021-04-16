package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/org"
)

func (c *Commands) AddMailText(ctx context.Context, resourceOwner string, mailText *domain.MailText) (*domain.MailText, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-MFiig", "Errors.ResourceOwnerMissing")
	}
	if !mailText.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-4778u", "Errors.Org.MailText.Invalid")
	}
	addedPolicy := NewOrgMailTextWriteModel(resourceOwner, mailText.MailTextType, mailText.Language)
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "Org-9kufs", "Errors.Org.MailText.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&addedPolicy.MailTextWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(
		ctx,
		org.NewMailTextAddedEvent(
			ctx,
			orgAgg,
			mailText.MailTextType,
			mailText.Language,
			mailText.Title,
			mailText.PreHeader,
			mailText.Subject,
			mailText.Greeting,
			mailText.Text,
			mailText.ButtonText))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToMailText(&addedPolicy.MailTextWriteModel), nil
}

func (c *Commands) ChangeMailText(ctx context.Context, resourceOwner string, mailText *domain.MailText) (*domain.MailText, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-NFus3", "Errors.ResourceOwnerMissing")
	}
	if !mailText.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-3m9fs", "Errors.Org.MailText.Invalid")
	}
	existingPolicy := NewOrgMailTextWriteModel(resourceOwner, mailText.MailTextType, mailText.Language)
	err := c.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "Org-3n8fM", "Errors.Org.MailText.NotFound")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.MailTextWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(
		ctx,
		orgAgg,
		mailText.MailTextType,
		mailText.Language,
		mailText.Title,
		mailText.PreHeader,
		mailText.Subject,
		mailText.Greeting,
		mailText.Text,
		mailText.ButtonText)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-2n9fs", "Errors.Org.MailText.NotChanged")
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToMailText(&existingPolicy.MailTextWriteModel), nil
}

func (c *Commands) RemoveMailText(ctx context.Context, resourceOwner, mailTextType, language string) error {
	if resourceOwner == "" {
		return caos_errs.ThrowInvalidArgument(nil, "Org-2N7fd", "Errors.ResourceOwnerMissing")
	}
	if mailTextType == "" || language == "" {
		return caos_errs.ThrowInvalidArgument(nil, "Org-N8fsf", "Errors.Org.MailText.Invalid")
	}
	existingPolicy := NewOrgMailTextWriteModel(resourceOwner, mailTextType, language)
	err := c.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return caos_errs.ThrowNotFound(nil, "Org-3b8Jf", "Errors.Org.MailText.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.WriteModel)
	_, err = c.eventstore.PushEvents(ctx, org.NewMailTextRemovedEvent(ctx, orgAgg, mailTextType, language))
	return err
}