package handler

import (
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type PhaseHandler struct {
	phaseSvc service.IPhaseService
}

func NewPhaseHandler(phaseSvc service.IPhaseService) *PhaseHandler {
	return &PhaseHandler{phaseSvc: phaseSvc}
}

// GetProjectWithPhases returns a project with all phases and their classes (public)
// GET /api/projects/:project_id
func (h *PhaseHandler) GetProjectWithPhases(c *fiber.Ctx) error {
	projectID := c.Params("project_id")
	if projectID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Project ID is required"})
	}

	project, err := h.phaseSvc.GetProjectWithPhases(c.UserContext(), projectID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Project not found"})
	}

	return c.JSON(project)
}

// ListPhasesByProject lists all phases under a project (public, no classes detail)
// GET /api/projects/:project_id/phases
func (h *PhaseHandler) ListPhasesByProject(c *fiber.Ctx) error {
	projectID := c.Params("project_id")
	if projectID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Project ID is required"})
	}

	phases, err := h.phaseSvc.ListPhasesByProject(c.UserContext(), projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": phases})
}

// GetPhase returns a single phase with its classes (public)
// GET /api/projects/:project_id/phases/:phase_id
func (h *PhaseHandler) GetPhase(c *fiber.Ctx) error {
	phaseID := c.Params("phase_id")
	if phaseID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Phase ID is required"})
	}

	phase, err := h.phaseSvc.GetPhase(c.UserContext(), phaseID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Phase not found"})
	}

	return c.JSON(phase)
}

// CreatePhase creates a new phase under a project (admin only)
// POST /api/admin/projects/:project_id/phases
func (h *PhaseHandler) CreatePhase(c *fiber.Ctx) error {
	adminID := c.Locals("userID").(string)
	projectID := c.Params("project_id")
	if projectID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Project ID is required"})
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		OrderIndex  int    `json:"order_index"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if input.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Title is required"})
	}

	phase, err := h.phaseSvc.CreatePhase(c.UserContext(), adminID, projectID, input.Title, input.Description, input.OrderIndex)
	if err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(phase)
}

// UpdatePhase updates phase metadata (admin only)
// PUT /api/admin/phases/:phase_id
func (h *PhaseHandler) UpdatePhase(c *fiber.Ctx) error {
	adminID := c.Locals("userID").(string)
	phaseID := c.Params("phase_id")
	if phaseID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Phase ID is required"})
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		OrderIndex  int    `json:"order_index"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	phase, err := h.phaseSvc.UpdatePhase(c.UserContext(), adminID, phaseID, input.Title, input.Description, input.OrderIndex)
	if err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(phase)
}

// DeletePhase deletes a phase and cascades its classes (admin only)
// DELETE /api/admin/phases/:phase_id
func (h *PhaseHandler) DeletePhase(c *fiber.Ctx) error {
	adminID := c.Locals("userID").(string)
	phaseID := c.Params("phase_id")
	if phaseID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Phase ID is required"})
	}

	if err := h.phaseSvc.DeletePhase(c.UserContext(), adminID, phaseID); err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Phase deleted"})
}
