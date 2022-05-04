package viewctx

type ViewContext struct {
  ScreenHeight int
  ScreenWidth  int
}

func NewViewContext() ViewContext {
  return ViewContext{}
}
