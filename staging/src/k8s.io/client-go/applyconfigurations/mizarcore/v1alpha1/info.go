/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// InfoApplyConfiguration represents an declarative configuration of the Info type for use
// with apply.
type InfoApplyConfiguration struct {
	Major        *string `json:"major,omitempty"`
	Minor        *string `json:"minor,omitempty"`
	GitVersion   *string `json:"gitVersion,omitempty"`
	GitCommit    *string `json:"gitCommit,omitempty"`
	GitTreeState *string `json:"gitTreeState,omitempty"`
	BuildDate    *string `json:"buildDate,omitempty"`
	GoVersion    *string `json:"goVersion,omitempty"`
	Compiler     *string `json:"compiler,omitempty"`
	Platform     *string `json:"platform,omitempty"`
}

// InfoApplyConfiguration constructs an declarative configuration of the Info type for use with
// apply.
func Info() *InfoApplyConfiguration {
	return &InfoApplyConfiguration{}
}

// WithMajor sets the Major field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Major field is set to the value of the last call.
func (b *InfoApplyConfiguration) WithMajor(value string) *InfoApplyConfiguration {
	b.Major = &value
	return b
}

// WithMinor sets the Minor field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Minor field is set to the value of the last call.
func (b *InfoApplyConfiguration) WithMinor(value string) *InfoApplyConfiguration {
	b.Minor = &value
	return b
}

// WithGitVersion sets the GitVersion field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the GitVersion field is set to the value of the last call.
func (b *InfoApplyConfiguration) WithGitVersion(value string) *InfoApplyConfiguration {
	b.GitVersion = &value
	return b
}

// WithGitCommit sets the GitCommit field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the GitCommit field is set to the value of the last call.
func (b *InfoApplyConfiguration) WithGitCommit(value string) *InfoApplyConfiguration {
	b.GitCommit = &value
	return b
}

// WithGitTreeState sets the GitTreeState field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the GitTreeState field is set to the value of the last call.
func (b *InfoApplyConfiguration) WithGitTreeState(value string) *InfoApplyConfiguration {
	b.GitTreeState = &value
	return b
}

// WithBuildDate sets the BuildDate field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the BuildDate field is set to the value of the last call.
func (b *InfoApplyConfiguration) WithBuildDate(value string) *InfoApplyConfiguration {
	b.BuildDate = &value
	return b
}

// WithGoVersion sets the GoVersion field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the GoVersion field is set to the value of the last call.
func (b *InfoApplyConfiguration) WithGoVersion(value string) *InfoApplyConfiguration {
	b.GoVersion = &value
	return b
}

// WithCompiler sets the Compiler field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Compiler field is set to the value of the last call.
func (b *InfoApplyConfiguration) WithCompiler(value string) *InfoApplyConfiguration {
	b.Compiler = &value
	return b
}

// WithPlatform sets the Platform field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Platform field is set to the value of the last call.
func (b *InfoApplyConfiguration) WithPlatform(value string) *InfoApplyConfiguration {
	b.Platform = &value
	return b
}
