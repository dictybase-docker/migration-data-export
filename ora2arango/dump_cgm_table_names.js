load("https://raw.githubusercontent.com/lodash/lodash/4.17.15-npm/core.min.js")
var FileWriter = Java.type("java.io.FileWriter")
var FileUtils = Java.type("oracle.dbtools.common.utils.FileUtils");
var cwd = FileUtils.getCWD(ctx);
var stmt = 'SELECT table_name name FROM all_tables WHERE owner = :name AND tablespace_name = :name '

var cgmfw = new FileWriter(cwd+"/cgm_ddb_tables.txt")
var ret2 = util.executeReturnList(stmt,{name:"CGM_DDB"}); 
_.forEach(ret2,function(row) { cgmfw.write(row.NAME + "\n")})
cgmfw.close()




