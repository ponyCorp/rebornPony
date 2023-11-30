package cmdhandler

func (h *CmdHandler) Mount() {
	h.Handle("help", "help", h.Help)
}
