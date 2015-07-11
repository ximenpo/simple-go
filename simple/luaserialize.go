package simple

import (
	"errors"
	. "github.com/ximenpo/simple-luago/lua"
)

func LuaSerialize(value interface{}) (result string, err error) {
	vm := LuaVM{}
	vm.Start()
	defer vm.Stop()

	vm.OpenStdLibs()

	if err = vm.RunString(LuaSerialize_SourceCode); err != nil {
		return
	}

	if !vm.SetObject("value", value, false) {
		return "", errors.New("set input object failed")
	}

	var ok bool
	var ref Lua_Ref
	if ref, ok = vm.Ref("value"); !ok {
		return "", errors.New("check input object failed")
	}

	err = vm.Invoke(&result, "serialize", &ref)
	return
}

const (
	LuaSerialize_SourceCode = `
serialize_pairsByKeys	= function(t, f)
    local a = {}
    for n in pairs(t) do table.insert(a, n) end
    table.sort(a, f)

    local i = 0                 -- iterator variable
    local iter = function ()    -- iterator function
        i = i + 1
        if a[i] == nil then return nil
        else return a[i], t[a[i]]
        end
    end
    return iter
end

serialize		= function (o, pre_tabs)
    local	code	= ""
    local	t		= pre_tabs or ''
    if type(o) == "number" then
        code = code .. o
    elseif type(o) == "boolean" then
        code = code .. tostring(o)
    elseif type(o) == "string" then
        code = code .. string.format("%q", o)
    elseif type(o) == "table" then
        code = code .. "\n"..t.."{\n"
        for k,v in serialize_pairsByKeys(o) do
            if k ~= '_G' and (type(v) == 'number' or type(v) == 'string' or type(v) == 'table') then
                if type(k) == 'number' then
                    code = code .. t.."    ["..k.."]\t= "
                    code = code .. serialize(v, t..'    ')
                    code = code .. ',\n'
                else
                    code = code .. t.."    "..k.."\t= "
                    code = code .. serialize(v, t..'    ')
                    code = code .. ',\n'
                end
            end
        end
        code = code .. t.."}"
    elseif type(o) == "function" then
        code = code .. t.."nil"
    else
        error("cannot serialize a " .. type(o))
    end
    return code
end
`
)
