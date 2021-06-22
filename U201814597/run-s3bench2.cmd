@rem -accessKey        Access Key
@rem -accessSecret     Secret Key
@rem -bucket=loadgen   Bucket for holding all test objects.
@rem -endpoint=http://127.0.0.1:9000 Endpoint URL of object storage service being tested.
@rem -numClients=8     Simulate 8 clients running concurrently.
@rem -numSamples=256   Test with 256 objects.
@rem -objectNamePrefix=loadgen Name prefix of test objects.
@rem -objectSize=1024          Size of test objects(byte)
@rem -verbose          Print latency for every request.

@echo off
@echo.>result2.txt
set /a numClients = 0
set /a objectSize = 1024
setlocal enabledelayedexpansion
for /l %%i in (0,1,19) do ( 
	set /a numClients +=1	
	s3bench.exe ^
		-accessKey=hust ^
		-accessSecret=hust_obs ^
		-bucket=loadgen ^
		-endpoint=http://127.0.0.1:9000 ^
		-numClients=!numClients! ^
		-numSamples=256 ^
		-objectNamePrefix=loadgen ^
		-objectSize=!objectSize! >> result2.txt
		echo !numClients!)



