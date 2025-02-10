var FileUtils = Java.type("oracle.dbtools.common.utils.FileUtils")
var Collectors = Java.type("java.util.stream.Collectors")
var Files = Java.type("java.nio.file.Files")
var Paths = Java.type("java.nio.file.Paths")
var BufferedReader = Java.type("java.io.BufferedReader")
var StandardOpenOption = Java.type("java.nio.file.StandardOpenOption")

ret = util.executeReturnList(stmt, {})
ret.forEach(function(row){
	ctx.write(row.FEATURE_ID)
	reader = row.RESIDUES.getCharacterStream()
	content = new BufferedReader(reader)
        	 .lines()
        	 .collect(Collectors.joining("\n"))
	Files.writeString(
		outputPath,
		String.join(",",row.FEATURE_ID,row.FEATURE_TYPE,content), 
		StandardOpenOption.APPEND
	)
	Files.writeString(
		outputPath,
		"\n",
		StandardOpenOption.APPEND
	)
	reader.close()
})
