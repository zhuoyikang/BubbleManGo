#! /usr/bin/awk -f
# 获取到所有的数据

BEGIN {
    LoadPath()
    LoadParam()
    LoadTypeAlias()
    ApiCount = 1
    PayLoadCount = 1
    LoadAPI()
    LoadPayLoad()
    # OutPutAPI()
    # OutPutPayLoad()
    OutPutBz()
}


# 从命令行参数获取中间文件和生成文件路径.
# PktGen.awk MiddlePath GenPath
function LoadPath() {
    middleFilePath = ARGV[1]
    genFilePath = ARGV[2]
}

function LoadParam() {
    file = middleFilePath"/name.txt"
    while( getline line < file ) { #由指定的文件中读取测验数据
        if (line ~ /^package*/  ) {
            split(line,a," ")
            packageParam=a[2]
        }

    }
}

# 类型别名.
function LoadTypeAlias() {
    TypeMap["byte"] = "byte"
    TypeMap["uint32"] = "uint32"
    TypeMap["int32"] = "int32"

    TypeMap["uint16"] = "uint16"
    TypeMap["int16"] = "int16"

    TypeMap["string"] = "string"
    TypeMap["int"] = "int32"
    TypeMap["integer"] = "int32"
}


function FindType(Type) {
    if (TypeMap[Type] != "") {
        return TypeMap[Type]
    }else {
        return Type
    }
}


# 加载api.txt
function LoadAPI() {
    file = middleFilePath"/api.txt"
    while( getline line < file ) { #由指定的文件中读取测验数据
        if (line ~ /^#.*/ || line ~ /^\s*$/ ) {
            if (ApiList[ApiCount,"packet_type"] != "") {
                ApiCount+=1
            }
            continue
        }
        split(line,a,":")
        ApiList[ApiCount,a[1]] = a[2]
    }
}


# 加载payload.txt
function LoadPayLoad() {
    file = middleFilePath"/proto.txt"
    while( getline line <  file) { #由指定的文件中读取测验数据
        if (line ~ /^#.*/ || line ~ /^\s*$/ || line == "===" ) {
            continue
        }
        if(match(line,/^[^=].+=/) > 0 ) {
            name=substr(line,0,length(line)-1)
            PayLoadList[name,"count"] = 0
            PayLoadNames[PayLoadCount] = name
            PayLoadCount+=1
        } else {
            split(line,a," ")
            if (a[2] == "array") {
                fc = PayLoadList[name,"count"]+1
                PayLoadList[name,fc,"name"]=a[1]
                PayLoadList[name,fc,"type"]="array"
                PayLoadList[name,fc,"addtion"]=FindType(a[3])
                PayLoadList[name,"count"]= fc
            } else {
                fc = PayLoadList[name,"count"]+1
                PayLoadList[name,fc,"name"]=a[1]
                PayLoadList[name,fc,"type"]=FindType(a[2])
                PayLoadList[name,"count"]= fc
            }
        }
    }
}

# 数据PayLoad
function OutPutPayLoad() {
    for (i = 1; i< PayLoadCount; i++) {
        name = PayLoadNames[i]
        print i,name
    }
}

# 输出go的struct结构.
function OutPutPayLoadStruct(Name) {
    count=PayLoadList[Name,"count"]
    printf("type %s struct {\n", Name) > genFile
    for(StructI=1; StructI<= count; StructI++) {
        name = PayLoadList[Name,StructI,"name"]
        type = PayLoadList[Name,StructI,"type"]
        addtion = PayLoadList[Name,StructI,"addtion"]
        if(type == "array") {
            if(TypeMap[addtion]=="") {
                addtion="*"addtion
            }
            printf("\t%s []%s\n", name,addtion) > genFile
        } else {
            if(TypeMap[type]=="") {
                type="*"type
            }
            printf("\t%s %s\n", name,type) > genFile
        }
    }
    print "}\n" > genFile
}

# 输出所有的PayLoadStruct.
function OutPutAllPayLoadStruct() {
    for(PayLoadi =1; PayLoadi< PayLoadCount; PayLoadi ++) {
        name = PayLoadNames[PayLoadi]
        OutPutPayLoadStruct(name)
    }
}

# 输出所有的API常量定义
function OutPutAllApiConst() {
    print "const (" > genFile
    for (APIi = 1; APIi <= ApiCount; APIi++) {
        printf("\tBZ_%s = %s\n", toupper(ApiList[APIi, "name"]),
               ApiList[APIi, "packet_type"]) > genFile
    }
    print ")\n" > genFile
}


# 输出所有的MapHandler
function OutPutMapHandler(Name) {
    printf("func MakeBz%sHandler() BzHandlerMAP {\n",Name)  > genFile
    print "\tProtocalHandler := BzHandlerMAP{" > genFile
    for(APIi=1; APIi <=ApiCount; APIi++) {
        if (ApiList[APIi, "name"] ~ /Req$/) {
            printf("\t\tBZ_%s: Bz%s,\n", toupper(ApiList[APIi, "name"]),
                   ApiList[APIi, "name"]) > genFile
        }
    }
    print "\t}\n\treturn ProtocalHandler\n}\n" > genFile
}


function OutPutErrCheck(prefix) {
    print prefix, "\tif err != nil {"  > genFile
    print prefix, "\t\treturn"  > genFile
    print prefix, "\t}"  > genFile
}

# 输出序列化:read
function OutPutSerializeRead(Name) {
    count=PayLoadList[Name,"count"]
    printf("func BzRead%s(datai []byte) (data []byte, ret *%s, err error) {\n",
           Name,Name) > genFile
    printf("\tdata = datai\n") > genFile
    printf("\tret = &%s{}\n",Name) > genFile
    for(Filedi = 1 ;Filedi<= count; Filedi++) {
        type = PayLoadList[Name,Filedi,"type"]
        addtion = PayLoadList[Name,Filedi,"addtion"]
        name = PayLoadList[Name,Filedi,"name"]

        if (type == "array") {

            if(TypeMap[addtion]=="") {
                addtion1="*"addtion
            }else{
                addtion1=addtion
            }
            printf("\tvar %s_v %s\n", name, addtion1) > genFile
            printf("\tdata, %s_size, err := BzReaduint16(data)\n", name) > genFile
            printf("\tfor i := 0; i < int(%s_size); i++ {\n", name)  > genFile
            printf("\t\tdata, %s_v, err = BzRead%s(data)\n",name,addtion) > genFile
            OutPutErrCheck("\t")
            printf("\t\tret.%s = append(ret.%s, %s_v)\n",name,name,name) > genFile
            printf("\t}\n") > genFile

        } else {
            printf("\tdata, ret.%s, err = BzRead%s(data)\n",name, type ) > genFile
        }
        OutPutErrCheck()
    }
    print "\treturn" > genFile
    printf("}\n") > genFile
}


# 输出序列化:write
function OutPutSerializeWrite(Name) {
    count=PayLoadList[Name,"count"]
    printf("func BzWrite%s(datai []byte, ret *%s) (data []byte, err error) {\n",
           Name,Name) > genFile
    printf("\tdata = datai\n") > genFile
    for(Filedi = 1 ;Filedi<= count; Filedi++) {
        type = PayLoadList[Name,Filedi,"type"]
        addtion = PayLoadList[Name,Filedi,"addtion"]
        name = PayLoadList[Name,Filedi,"name"]
        if (type == "array") {
            printf("\tdata, err = BzWriteuint16(data, uint16(len(ret.%s)))\n",
                   name) > genFile
            printf("\tfor _, %s_v := range ret.%s {\n", name, name)  > genFile
            printf("\t\tdata, err = BzWrite%s(data, %s_v)\n",addtion,name) > genFile
            printf("\t}\n") > genFile

        } else {
            printf("\tdata, err = BzWrite%s(data, ret.%s)\n",type, name) > genFile
        }
    }
    print "\treturn" > genFile
    printf("}\n") > genFile
}


# 输出序列化函数.
function OutPutSerialize(Name) {
    OutPutSerializeRead(Name)
    OutPutSerializeWrite(Name)
}

# 输出所有的序列化函数.
function OutPutAllSerialize() {
    for (i = 1; i< PayLoadCount; i++) {
        Name = PayLoadNames[i]
        OutPutSerialize(Name)
    }
}

# 输出BzGo文件.
function OutPutBz() {
    genFile = genFilePath "/demo.go"
    print "package " packageParam "\n" > genFile
    print "import (. \"agent\")\n"  > genFile
    OutPutAllApiConst()
    OutPutAllPayLoadStruct()
#    OutPutMapHandler("Gs")
    OutPutAllSerialize()
}

{
}
