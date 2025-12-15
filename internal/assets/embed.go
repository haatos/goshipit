package assets

import "embed"

//go:embed public/*
var PublicFS embed.FS

//go:embed content/*
var ContentFS embed.FS

//go:embed generated/*
var GeneratedFS embed.FS
