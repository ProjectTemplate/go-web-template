package response

import (
	"bytes"
	"errors"
	"go-web-template/base/common/utils"
)

const (
	// CodeLengthProjectModule 项目模块部分编码长度
	codeLengthProjectModule = 3
	// CodeLengthProject 项目部分编码长度
	codeLengthProject = 3
	// CodeLengthDepartment 部门部分编码长度
	codeLengthDepartment = 3
	// CodeLengthCompany 公司部分长度
	codeLengthCompany = 3
	// CodeLength 错误码长度
	codeLength = 4

	// totalLength 总长度
	totalLength        = codeLengthProjectModule + codeLengthProject + codeLengthDepartment + codeLengthCompany + codeLength
	codeStart          = totalLength - codeLength
	projectModuleStart = codeStart - codeLengthProjectModule
	projectStart       = projectModuleStart - codeLengthProject
	departmentStart    = projectStart - codeLengthProject
	companyStart       = departmentStart - codeLengthCompany
)

const (
	// CodeCommon 通用错误码
	CodeCommon = "0000"
	//CodeInternalError 内部错误
	CodeInternalError = "0001"

	//CodeParamError 参数错误
	CodeParamError = "0002"
)

var codeMessageMap = map[string]string{
	CodeCommon:        "通用错误",
	CodeInternalError: "内部错误",
	CodeParamError:    "参数错误",
}

var (
	SuccessCode = ""

	// AdminCode 管理员项目错误码
	AdminCode              = NewCode("101", "001", "001", "001").WithCode(CodeCommon)
	AdminInternalErrorCode = AdminCode.WithCode(CodeInternalError)
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

func (c Code) Encode() string {
	return Encode(c)
}

func Encode(c Code) string {
	bufLength := totalLength
	buffer := bytes.NewBuffer(make([]byte, 0, bufLength))

	//数据格式校验
	err := check(c)
	utils.PanicAndPrintIfNotNil(err)

	c.ProjectModule = utils.FillZero(c.ProjectModule, codeLengthProjectModule)
	c.Project = utils.FillZero(c.Project, codeLengthProject)
	c.Department = utils.FillZero(c.Department, codeLengthDepartment)
	c.Company = utils.FillZero(c.Company, codeLengthCompany)
	c.Code = utils.FillZero(c.Code, codeLengthCompany)

	buffer.WriteString(c.Company)
	buffer.WriteString(c.Department)
	buffer.WriteString(c.Project)
	buffer.WriteString(c.ProjectModule)
	buffer.WriteString(c.Code)

	return buffer.String()
}

func Decode(code string) Code {
	result := Code{}

	if len(code) != totalLength {
		utils.PanicAndPrintIfNotNil(errors.New("code length is illegal"))
	}

	result.Code = code[codeStart:totalLength]
	result.ProjectModule = code[projectModuleStart:codeStart]
	result.Project = code[projectStart:projectModuleStart]
	result.Department = code[departmentStart:projectStart]
	result.Company = code[companyStart:departmentStart]

	return result
}

func check(c Code) error {
	//数据格式校验
	if len(c.ProjectModule) > codeLengthProjectModule {
		return errors.New("project module length is illegal")
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

	if len(c.Code) > codeLength {
		return errors.New("code length is illegal")
	}

	return nil
}
