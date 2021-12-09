# Compile
javac ButtonDemo.java -d .

jar cmf ButtonDemo.mf ButtonDemo.jar components/

# Run Compiled
java -jar ButtonDemo.jar

# Run Source (clean any binary classes or jar)
java ButtonDemo.java
