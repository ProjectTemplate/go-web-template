package response

import (
	"bytes"
	"errors"
	"go-web-template/base/common/utils"
)

const (
	// CodeLength 错误码长度
	codeLengthCode = 4
	// CodeLengthProjectModule 项目模块部分编码长度
	codeLengthProjectModule = 3
	// CodeLengthProject 项目部分编码长度
	codeLengthProject = 3
	// CodeLengthDepartment 部门部分编码长度
	codeLengthDepartment = 3
	// CodeLengthCompany 公司部分长度
	codeLengthCompany = 3

	// totalLength 总长度
	totalLength        = codeLengthProjectModule + codeLengthProject + codeLengthDepartment + codeLengthCompany + codeLengthCode
	codeStart          = totalLength - codeLengthCode
	projectModuleStart = codeStart - codeLengthProjectModule
	projectStart       = projectModuleStart - codeLengthProject
	departmentStart    = projectStart - codeLengthProject
	companyStart       = departmentStart - codeLengthCompany
)

const (
	// CodeCommon 通用错误码
	CodeCommon = "0000"
	// CodeInternalError 内部错误
	CodeInternalError = "0001"

	// CodeParamError 参数错误
	CodeParamError = "1001"

	// CodeDBError 数据库错误
	CodeDBError = "1101"

	// CodeKafKaError kafka错误
	CodeKafKaError = "1201"

	// CodeRedisError redis错误
	CodeRedisError = "1301"

	// CodeNacosError nacos错误
	CodeNacosError = "1401"

	// CodeESError ES错误
	CodeESError = "1501"
)

var codeMessageMap = map[string]string{
	CodeCommon:        "通用错误",
	CodeInternalError: "内部错误",
	CodeParamError:    "参数错误",
	CodeDBError:       "数据库错误",
	CodeKafKaError:    "kafka错误",
	CodeRedisError:    "redis错误",
	CodeNacosError:    "nacos错误",
	CodeESError:       "ES错误",
}

var (
	SuccessCode = ""

	// AdminCode 管理员项目错误码
	AdminCode              = NewCode("101", "001", "001", "001").WithCode(CodeCommon)
	AdminCodeInternalError = AdminCode.WithCode(CodeInternalError)
	AdminCodeParamError    = AdminCode.WithCode(CodeParamError)
)

// Code 编码
type Code struct {
	// ProjectModule  项目模块
	ProjectModule string
	// Project 项目
	Project string
	// Department 部门
	Department string
	// Company 公司
	Company string
	// Code 错误码
	Code string
}

func NewCode(company, department, project, projectModule string) Code {
	code := Code{
		ProjectModule: projectModule,
		Project:       project,
		Department:    department,
		Company:       company,
	}

	//校验数据是否合法
	err := check(code)
	utils.PanicAndPrintIfNotNil(err)

	return code
}

// WithCode 设置错误码
func (c Code) WithCode(code string) Code {
	result := Code{
		Department:    c.Department,
		Company:       c.Company,
		ProjectModule: c.ProjectModule,
		Project:       c.Project,
		Code:          code,
	}
	err := check(result)
	utils.PanicAndPrintIfNotNil(err)
	return result
}

// Encode 编码
func (c Code) Encode() string {
	return Encode(c)
}

// Encode 编码
func Encode(c Code) string {
	bufLength := totalLength
	buffer := bytes.NewBuffer(make([]byte, 0, bufLength))

	//数据格式校验
	err := check(c)
	utils.PanicAndPrintIfNotNil(err)

	c.ProjectModule = utils.FillZeroToNumberString(c.ProjectModule, codeLengthProjectModule)
	c.Project = utils.FillZeroToNumberString(c.Project, codeLengthProject)
	c.Department = utils.FillZeroToNumberString(c.Department, codeLengthDepartment)
	c.Company = utils.FillZeroToNumberString(c.Company, codeLengthCompany)
	c.Code = utils.FillZeroToNumberString(c.Code, codeLengthCode)

	buffer.WriteString(c.Company)
	buffer.WriteString(c.Department)
	buffer.WriteString(c.Project)
	buffer.WriteString(c.ProjectModule)
	buffer.WriteString(c.Code)

	return buffer.String()
}

// Decode 解码
func Decode(code string) Code {
	result := Code{}

	if len(code) != totalLength {
		panic(errors.New("code length is illegal"))
	}

	result.Code = code[codeStart:totalLength]
	result.ProjectModule = code[projectModuleStart:codeStart]
	result.Project = code[projectStart:projectModuleStart]
	result.Department = code[departmentStart:projectStart]
	result.Company = code[companyStart:departmentStart]

	return result
}

// check 数据格式校验
func check(c Code) error {

	if len(c.ProjectModule) > codeLengthProjectModule {
		return errors.New("project model length is illegal")
	}

	if len(c.Project) > codeLengthProject {
		return errors.New("project length is illegal")
	}

	if len(c.Department) > codeLengthDepartment {
		return errors.New("department length is illegal")
	}

	if len(c.Company) > codeLengthCompany {
		return errors.New("company length is illegal")
	}

	if len(c.Code) > codeLengthCode {
		return errors.New("code length is illegal")
	}

	return nil
}
