package custom_errors

import "errors"

var ErrNotFound = errors.New("not found")

var ErrAlreadyExists = errors.New("already exists")

var ErrDatabaseFailure = errors.New("database failure")

var ErrExternalServiceFailure = errors.New("external service failure")

// Ошибки, связанные с сущностями
var ErrEntityNotFound = errors.New("entity not found")
var ErrEntityAlreadyExists = errors.New("entity already exists")

// Ошибки, связанные с инфраструктурой
var ErrRepoQueryFailed = errors.New("repository query failed")
var ErrRepoScanFailed = errors.New("repository scan failed")

var ErrPermissionDenied = errors.New("permission denied")
