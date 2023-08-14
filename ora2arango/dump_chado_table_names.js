load("https://raw.githubusercontent.com/lodash/lodash/4.17.15-npm/core.min.js")
var FileWriter = Java.type("java.io.FileWriter")
var FileUtils = Java.type("oracle.dbtools.common.utils.FileUtils");
var cwd = FileUtils.getCWD(ctx);
var stmt = 'SELECT table_name name FROM all_tables WHERE owner = :name AND tablespace_name = :name '

var chadofw = new FileWriter(cwd+"/chado_tables.txt")
var ret = util.executeReturnList(stmt,{name:"CGM_CHADO"}); 
_.forEach(ret,function(row) { chadofw.write(row.NAME + "\n")})
chadofw.close()
