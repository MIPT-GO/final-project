package custom_errors

import "errors"

var ErrNotFound = errors.New("not found")

var ErrAlreadyExists = errors.New("already exists")

var ErrDatabaseFailure = errors.New("database failure")

var ErrExternalServiceFailure = errors.New("external service failure")

var ErrEntityNotFound = errors.New("entity not found")
var ErrEntityAlreadyExists = errors.New("entity already exists")

var ErrRepoQueryFailed = errors.New("repository query failed")
var ErrRepoScanFailed = errors.New("repository scan failed")

var ErrPermissionDenied = errors.New("permission denied")
