// Code generated by ogen, DO NOT EDIT.

package generated

import (
	"net/http"
	"net/url"

	"github.com/go-faster/errors"

	"github.com/ogen-go/ogen/conv"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/uri"
	"github.com/ogen-go/ogen/validate"
)

// GetMeasurementsByStationParams is parameters of getMeasurementsByStation operation.
type GetMeasurementsByStationParams struct {
	ID int32
}

func unpackGetMeasurementsByStationParams(packed middleware.Parameters) (params GetMeasurementsByStationParams) {
	{
		key := middleware.ParameterKey{
			Name: "id",
			In:   "path",
		}
		params.ID = packed[key].(int32)
	}
	return params
}

func decodeGetMeasurementsByStationParams(args [1]string, argsEscaped bool, r *http.Request) (params GetMeasurementsByStationParams, _ error) {
	// Decode path: id.
	if err := func() error {
		param := args[0]
		if argsEscaped {
			unescaped, err := url.PathUnescape(args[0])
			if err != nil {
				return errors.Wrap(err, "unescape path")
			}
			param = unescaped
		}
		if len(param) > 0 {
			d := uri.NewPathDecoder(uri.PathDecoderConfig{
				Param:   "id",
				Value:   param,
				Style:   uri.PathStyleSimple,
				Explode: false,
			})

			if err := func() error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToInt32(val)
				if err != nil {
					return err
				}

				params.ID = c
				return nil
			}(); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "id",
			In:   "path",
			Err:  err,
		}
	}
	return params, nil
}

// GetParametersByStationParams is parameters of getParametersByStation operation.
type GetParametersByStationParams struct {
	ID int32
}

func unpackGetParametersByStationParams(packed middleware.Parameters) (params GetParametersByStationParams) {
	{
		key := middleware.ParameterKey{
			Name: "id",
			In:   "path",
		}
		params.ID = packed[key].(int32)
	}
	return params
}

func decodeGetParametersByStationParams(args [1]string, argsEscaped bool, r *http.Request) (params GetParametersByStationParams, _ error) {
	// Decode path: id.
	if err := func() error {
		param := args[0]
		if argsEscaped {
			unescaped, err := url.PathUnescape(args[0])
			if err != nil {
				return errors.Wrap(err, "unescape path")
			}
			param = unescaped
		}
		if len(param) > 0 {
			d := uri.NewPathDecoder(uri.PathDecoderConfig{
				Param:   "id",
				Value:   param,
				Style:   uri.PathStyleSimple,
				Explode: false,
			})

			if err := func() error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToInt32(val)
				if err != nil {
					return err
				}

				params.ID = c
				return nil
			}(); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "id",
			In:   "path",
			Err:  err,
		}
	}
	return params, nil
}

// GetStationByIdParams is parameters of getStationById operation.
type GetStationByIdParams struct {
	ID int32
}

func unpackGetStationByIdParams(packed middleware.Parameters) (params GetStationByIdParams) {
	{
		key := middleware.ParameterKey{
			Name: "id",
			In:   "path",
		}
		params.ID = packed[key].(int32)
	}
	return params
}

func decodeGetStationByIdParams(args [1]string, argsEscaped bool, r *http.Request) (params GetStationByIdParams, _ error) {
	// Decode path: id.
	if err := func() error {
		param := args[0]
		if argsEscaped {
			unescaped, err := url.PathUnescape(args[0])
			if err != nil {
				return errors.Wrap(err, "unescape path")
			}
			param = unescaped
		}
		if len(param) > 0 {
			d := uri.NewPathDecoder(uri.PathDecoderConfig{
				Param:   "id",
				Value:   param,
				Style:   uri.PathStyleSimple,
				Explode: false,
			})

			if err := func() error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToInt32(val)
				if err != nil {
					return err
				}

				params.ID = c
				return nil
			}(); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "id",
			In:   "path",
			Err:  err,
		}
	}
	return params, nil
}
