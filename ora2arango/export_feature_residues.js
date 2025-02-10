FileUtils = Java.type("oracle.dbtools.common.utils.FileUtils")
Collectors = Java.type("java.util.stream.Collectors")
Files = Java.type("java.nio.file.Files")
Paths = Java.type("java.nio.file.Paths")
BufferedReader = Java.type("java.io.BufferedReader")
StandardOpenOption = Java.type("java.nio.file.StandardOpenOption")
String = Java.type("java.lang.String")

cwd = FileUtils.getCWD(ctx);
outputPath = Paths.get(cwd, "/","feature_clob.csv")
Files.writeString(outputPath, "", StandardOpenOption.CREATE, StandardOpenOption.TRUNCATE_EXISTING)

var stmt = `
SELECT 
    feature_id,
    cvterm.name as feature_type,
    residues
FROM feature 
JOIN cvterm ON cvterm.cvterm_id = feature.type_id
WHERE residues IS NOT NULL
`

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
