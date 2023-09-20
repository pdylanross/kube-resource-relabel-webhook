package mutation

import (
	"encoding/json"
	"fmt"

	"gomodules.xyz/jsonpatch/v3"

	"github.com/pdylanross/kube-resource-relabel-webhook/pkg/util"
	admissionv1 "k8s.io/api/admission/v1"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
)

var admissionScheme = runtime.NewScheme()
var admissionCodecs = serializer.NewCodecFactory(admissionScheme)
var universalDeserializer = clientsetscheme.Codecs.UniversalDeserializer()

func init() {
	util.ErrCheck(admissionv1.AddToScheme(admissionScheme))
	util.ErrCheck(admissionv1beta1.AddToScheme(admissionScheme))
}

type AdmissionReview interface {
	GetObject() (metaV1.Object, error)
	GetRawObject() []byte
	SetError(err error)
	ToSerializeable() interface{}
	SetPatches(patches []jsonpatch.Operation) error
	IsDryRun() bool
	GetStatus() int
}

func NewAdmissionReview(body []byte) (AdmissionReview, error) {
	review, _, err := admissionCodecs.UniversalDeserializer().Decode(body, nil, nil)

	if err != nil {
		return nil, err
	}

	switch r := review.(type) {
	case *admissionv1.AdmissionReview:
		review := admissionReviewV1{original: r}
		return &review, nil
	case *admissionv1beta1.AdmissionReview:
		review := admissionReviewV1Beta1{original: r}
		return &review, nil
	}

	return nil, fmt.Errorf("unhandled admission review type: %s", review.GetObjectKind().GroupVersionKind().String())
}

type admissionReviewV1 struct {
	original *admissionv1.AdmissionReview
	response *admissionv1.AdmissionResponse
}

func (a *admissionReviewV1) GetRawObject() []byte {
	return a.original.Request.Object.Raw
}

func (a *admissionReviewV1) GetStatus() int {
	if a.response == nil {
		return 200
	}

	if a.response.Allowed && len(a.response.Warnings) > 0 {
		return 299
	} else if a.response.Allowed {
		return 200
	} else {
		return int(a.response.Result.Code)
	}
}

func (a *admissionReviewV1) SetPatches(patches []jsonpatch.Operation) error {
	raw, err := json.Marshal(patches)
	if err != nil {
		return fmt.Errorf("cound not marshal json patch to json: %w", err)
	}

	if a.response == nil {
		a.response = &admissionv1.AdmissionResponse{UID: a.original.Request.UID, Allowed: true}
	}

	pt := admissionv1.PatchTypeJSONPatch

	a.response.Patch = raw
	a.response.PatchType = &pt

	return nil
}

func (a *admissionReviewV1) IsDryRun() bool {
	return *a.original.Request.DryRun
}

func (a *admissionReviewV1) GetObject() (metaV1.Object, error) {
	return decodeK8sObject(a.original.Request.Object.Raw)
}

func (a *admissionReviewV1) SetError(err error) {
	if a.response == nil {
		a.response = &admissionv1.AdmissionResponse{UID: a.original.Request.UID}
	}

	a.response.Allowed = false
	a.response.Result = &metaV1.Status{
		Message: fmt.Sprintf("unhandled error with kube-resource-relabel-webhook: %s", err.Error()),
		Code:    500,
	}
}

func (a *admissionReviewV1) ToSerializeable() interface{} {
	if a.response == nil {
		a.response = &admissionv1.AdmissionResponse{UID: a.original.Request.UID, Allowed: true}
	}

	return admissionv1.AdmissionReview{Response: a.response, TypeMeta: a.original.TypeMeta}
}

type admissionReviewV1Beta1 struct {
	original *admissionv1beta1.AdmissionReview
	response *admissionv1beta1.AdmissionResponse
}

func (a *admissionReviewV1Beta1) GetRawObject() []byte {
	return a.original.Request.Object.Raw
}

func (a *admissionReviewV1Beta1) GetStatus() int {
	if a.response == nil {
		return 200
	}

	if a.response.Allowed && len(a.response.Warnings) > 0 {
		return 299
	} else if a.response.Allowed {
		return 200
	} else {
		return int(a.response.Result.Code)
	}
}

func (a *admissionReviewV1Beta1) SetPatches(patches []jsonpatch.Operation) error {
	raw, err := json.Marshal(patches)
	if err != nil {
		return fmt.Errorf("cound not marshal json patch to json: %w", err)
	}

	if a.response == nil {
		a.response = &admissionv1beta1.AdmissionResponse{UID: a.original.Request.UID, Allowed: true}
	}

	pt := admissionv1beta1.PatchTypeJSONPatch

	a.response.Patch = raw
	a.response.PatchType = &pt

	return nil
}

func (a *admissionReviewV1Beta1) IsDryRun() bool {
	return *a.original.Request.DryRun
}

func (a *admissionReviewV1Beta1) GetObject() (metaV1.Object, error) {
	return decodeK8sObject(a.original.Request.Object.Raw)
}

func (a *admissionReviewV1Beta1) SetError(err error) {
	if a.response == nil {
		a.response = &admissionv1beta1.AdmissionResponse{UID: a.original.Request.UID}
	}

	a.response.Allowed = false
	a.response.Result = &metaV1.Status{
		Message: fmt.Sprintf("unhandled error with kube-resource-relabel-webhook: %s", err.Error()),
		Code:    500,
	}
}

func (a *admissionReviewV1Beta1) ToSerializeable() interface{} {
	if a.response == nil {
		a.response = &admissionv1beta1.AdmissionResponse{UID: a.original.Request.UID, Allowed: true}
	}

	return admissionv1beta1.AdmissionReview{Response: a.response, TypeMeta: a.original.TypeMeta}
}

func decodeK8sObject(raw []byte) (metaV1.Object, error) {
	des, _, err := universalDeserializer.Decode(raw, nil, nil)
	if err != nil {
		return nil, err
	}

	obj, ok := des.(metaV1.Object)
	if !ok {
		return nil, fmt.Errorf("could not type assert metav1.Object")
	}

	return obj, nil
}
