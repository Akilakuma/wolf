smart={}

smart.parser_string = function(para)

--     使用golang提供的method
    local inm = require("inmodule")
    local json_value = inm.gjson(para, "number")

    return json_value
end