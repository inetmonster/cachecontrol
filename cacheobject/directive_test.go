/**
 *  Copyright 2015 Paul Querna
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package cacheobject

import (
	"github.com/stretchr/testify/require"

	"fmt"
	"math"
	"testing"
)

func TestMaxAge(t *testing.T) {
	cd, err := ParseResponseCacheControl("")
	require.NoError(t, err)
	require.Equal(t, cd.MaxAge, -1)

	cd, err = ParseResponseCacheControl("max-age")
	require.Error(t, err)

	cd, err = ParseResponseCacheControl("max-age=20")
	require.NoError(t, err)
	require.Equal(t, cd.MaxAge, 20)

	cd, err = ParseResponseCacheControl("max-age=0")
	require.NoError(t, err)
	require.Equal(t, cd.MaxAge, 0)

	cd, err = ParseResponseCacheControl("max-age=-1")
	require.Error(t, err)
}

func TestSMaxAge(t *testing.T) {
	cd, err := ParseResponseCacheControl("")
	require.NoError(t, err)
	require.Equal(t, cd.SMaxAge, -1)

	cd, err = ParseResponseCacheControl("s-maxage")
	require.Error(t, err)

	cd, err = ParseResponseCacheControl("s-maxage=20")
	require.NoError(t, err)
	require.Equal(t, cd.SMaxAge, 20)

	cd, err = ParseResponseCacheControl("s-maxage=0")
	require.NoError(t, err)
	require.Equal(t, cd.SMaxAge, 0)

	cd, err = ParseResponseCacheControl("s-maxage=-1")
	require.Error(t, err)
}

func TestResNoCache(t *testing.T) {
	cd, err := ParseResponseCacheControl("")
	require.NoError(t, err)
	require.Equal(t, cd.SMaxAge, -1)

	cd, err = ParseResponseCacheControl("no-cache")
	require.NoError(t, err)
	require.Equal(t, cd.NoCachePresent, true)
	require.Equal(t, len(cd.NoCache), 0)

	cd, err = ParseResponseCacheControl("no-cache=MyThing")
	require.NoError(t, err)
	require.Equal(t, cd.NoCachePresent, true)
	require.Equal(t, len(cd.NoCache), 1)
}

func TestResSpaceOnly(t *testing.T) {
	cd, err := ParseResponseCacheControl(" ")
	require.NoError(t, err)
	require.Equal(t, cd.SMaxAge, -1)
}

func TestResTabOnly(t *testing.T) {
	cd, err := ParseResponseCacheControl("\t")
	require.NoError(t, err)
	require.Equal(t, cd.SMaxAge, -1)
}

func TestResPrivateExtensionQuoted(t *testing.T) {
	cd, err := ParseResponseCacheControl(`private="Set-Cookie,Request-Id" public`)
	require.NoError(t, err)
	require.Equal(t, cd.Public, true)
	require.Equal(t, cd.PrivatePresent, true)
	require.Equal(t, len(cd.Private), 2)
	require.Equal(t, len(cd.Extensions), 0)
	require.Equal(t, cd.Private["Set-Cookie"], true)
	require.Equal(t, cd.Private["Request-Id"], true)
}

func TestResPrivateExtension(t *testing.T) {
	cd, err := ParseResponseCacheControl(`private=Set-Cookie,Request-Id public`)
	require.NoError(t, err)
	require.Equal(t, cd.Public, true)
	require.Equal(t, cd.PrivatePresent, true)
	require.Equal(t, len(cd.Private), 2)
	require.Equal(t, len(cd.Extensions), 0)
	require.Equal(t, cd.Private["Set-Cookie"], true)
	require.Equal(t, cd.Private["Request-Id"], true)
}

func TestResMultipleNoCacheTabExtension(t *testing.T) {
	cd, err := ParseResponseCacheControl("no-cache " + "\t" + "no-cache=Mything aasdfdsfa")
	require.NoError(t, err)
	require.Equal(t, cd.NoCachePresent, true)
	require.Equal(t, len(cd.NoCache), 1)
	require.Equal(t, len(cd.Extensions), 1)
	require.Equal(t, cd.NoCache["Mything"], true)
}

func TestResExtensionsEmptyQuote(t *testing.T) {
	cd, err := ParseResponseCacheControl(`foo="" bar="hi"`)
	require.NoError(t, err)
	require.Equal(t, cd.SMaxAge, -1)
	require.Equal(t, len(cd.Extensions), 2)
	require.Contains(t, cd.Extensions, "bar=hi")
	require.Contains(t, cd.Extensions, "foo=")
}

func TestResQuoteMismatch(t *testing.T) {
	cd, err := ParseResponseCacheControl(`foo="`)
	require.Error(t, err)
	require.Nil(t, cd)
	require.Equal(t, err, ErrQuoteMismatch)
}

func TestResMustRevalidateNoArgs(t *testing.T) {
	cd, err := ParseResponseCacheControl(`must-revalidate=234`)
	require.Error(t, err)
	require.Nil(t, cd)
	require.Equal(t, err, ErrMustRevalidateNoArgs)
}

func TestResNoTransformNoArgs(t *testing.T) {
	cd, err := ParseResponseCacheControl(`no-transform="xxx"`)
	require.Error(t, err)
	require.Nil(t, cd)
	require.Equal(t, err, ErrNoTransformNoArgs)
}

func TestResNoStoreNoArgs(t *testing.T) {
	cd, err := ParseResponseCacheControl(`no-store=""`)
	require.Error(t, err)
	require.Nil(t, cd)
	require.Equal(t, err, ErrNoStoreNoArgs)
}

func TestResProxyRevalidateNoArgs(t *testing.T) {
	cd, err := ParseResponseCacheControl(`proxy-revalidate=23432`)
	require.Error(t, err)
	require.Nil(t, cd)
	require.Equal(t, err, ErrProxyRevalidateNoArgs)
}

func TestResPublicNoArgs(t *testing.T) {
	cd, err := ParseResponseCacheControl(`public=999Vary`)
	require.Error(t, err)
	require.Nil(t, cd)
	require.Equal(t, err, ErrPublicNoArgs)
}

func TestResMustRevalidate(t *testing.T) {
	cd, err := ParseResponseCacheControl(`must-revalidate`)
	require.NoError(t, err)
	require.NotNil(t, cd)
	require.Equal(t, cd.MustRevalidate, true)
}

func TestResNoTransform(t *testing.T) {
	cd, err := ParseResponseCacheControl(`no-transform`)
	require.NoError(t, err)
	require.NotNil(t, cd)
	require.Equal(t, cd.NoTransform, true)
}

func TestResNoStore(t *testing.T) {
	cd, err := ParseResponseCacheControl(`no-store`)
	require.NoError(t, err)
	require.NotNil(t, cd)
	require.Equal(t, cd.NoStore, true)
}

func TestResProxyRevalidate(t *testing.T) {
	cd, err := ParseResponseCacheControl(`proxy-revalidate`)
	require.NoError(t, err)
	require.NotNil(t, cd)
	require.Equal(t, cd.ProxyRevalidate, true)
}

func TestResPublic(t *testing.T) {
	cd, err := ParseResponseCacheControl(`public`)
	require.NoError(t, err)
	require.NotNil(t, cd)
	require.Equal(t, cd.Public, true)
}

func TestResPrivate(t *testing.T) {
	cd, err := ParseResponseCacheControl(`private`)
	require.NoError(t, err)
	require.NotNil(t, cd)
	require.Len(t, cd.Private, 0)
	require.Equal(t, cd.PrivatePresent, true)
}

func TestParseDeltaSecondsZero(t *testing.T) {
	ds, err := parseDeltaSeconds("0")
	require.NoError(t, err)
	require.Equal(t, ds, 0)
}

func TestParseDeltaSecondsLarge(t *testing.T) {
	ds, err := parseDeltaSeconds(fmt.Sprintf("%d", int64(math.MaxInt32)*2))
	require.NoError(t, err)
	require.Equal(t, ds, math.MaxInt32)
}

func TestParseDeltaSecondsVeryLarge(t *testing.T) {
	ds, err := parseDeltaSeconds(fmt.Sprintf("%d", math.MaxInt64))
	require.NoError(t, err)
	require.Equal(t, ds, math.MaxInt32)
}

func TestParseDeltaSecondsNegative(t *testing.T) {
	ds, err := parseDeltaSeconds("-60")
	require.Error(t, err)
	require.Equal(t, DeltaSeconds(-1), ds)
}
