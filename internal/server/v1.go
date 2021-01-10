package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sharmarajdaksh/yorpoll-api/internal/poll"
)

func (h *handler) getPoll(c *fiber.Ctx) error {
	pollID := c.Params("pollID")
	err := vldt.Var(pollID, "required,len=36")
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return fiber.ErrNotFound
	}

	p, err := h.dbc.GetPollByID(pollID)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return fiber.ErrInternalServerError
	}
	if p == nil {
		c.Status(fiber.StatusNotFound)
		return fiber.ErrNotFound
	}

	return c.Status(200).JSON(p)
}

func (h *handler) postPoll(c *fiber.Ctx) error {
	p := struct {
		Title       string   `json:"title" validate:"required,min=1"`
		Description string   `json:"description" validate:"required,min=1"`
		Expiry      int64    `json:"expiry" validate:"required,number,min=0"`
		Options     []string `json:"options" validate:"required,min=1,dive,min=1"`
	}{}

	if err := c.BodyParser(&p); err != nil {
		c.Status(fiber.StatusBadRequest)
		return fiber.ErrBadRequest
	}

	err := validate(p)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return err
	}

	opts := []poll.Option{}
	for _, o := range p.Options {
		opt := poll.NewOption(o)
		opts = append(opts, opt)
	}

	np := poll.New(p.Title, p.Description, opts, p.Expiry)

	err = h.dbc.SavePoll(&np)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusCreated).JSON(np)
}

func (h *handler) deletePoll(c *fiber.Ctx) error {
	pollID := c.Params("pollID")
	err := vldt.Var(pollID, idValidationString)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return fiber.ErrNotFound
	}

	d, err := h.dbc.DeletePollByID(pollID)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return fiber.ErrInternalServerError
	}
	if !d {
		c.Status(fiber.StatusNotFound)
		return fiber.ErrNotFound
	}

	return c.SendStatus(http.StatusAccepted)
}

func (h *handler) putVote(c *fiber.Ctx) error {
	// expected params
	pollID := c.Params("pollID")
	optionID := c.Params("optionID")
	err := vldt.Var(pollID, idValidationString)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return fiber.ErrNotFound
	}
	err = vldt.Var(optionID, idValidationString)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return fiber.ErrNotFound
	}

	a, err := h.dbc.AddPollVote(pollID, optionID)
	if err != nil {
		if err == poll.ErrPollExpired {
			c.Status(fiber.StatusLocked)
			return fiber.NewError(fiber.StatusLocked, idValidationString)
		}
		c.Status(fiber.StatusInternalServerError)
		return fiber.ErrInternalServerError
	}
	if !a {
		c.Status(fiber.StatusNotFound)
		return fiber.ErrNotFound
	}

	return c.SendStatus(http.StatusOK)
}
