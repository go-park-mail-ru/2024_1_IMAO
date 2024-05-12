vegeta attack -duration=5s -rate=1000 -targets=http://127.0.0.1:8008
echo "GET http://127.0.0.1:8008" | vegeta attack -name=50qps -rate=50 -duration=5s > results.50qps.bin
echo "GET http://localhost:8080/api/adverts/?count=20&startId=2000" | vegeta attack -name=50qps -rate=50 -duration=3s > results_50qps.csv
echo "GET http://localhost:8080/api/adverts/?count=20&startId=2000" | vegeta attack -name=100qps -rate=100 -duration=3s > results_100qps.csv
echo "GET http://localhost:8080/api/adverts/?count=20&startId=2000" | vegeta attack -name=50qps -rate=50 -duration=3s --output results.bin; 
vegeta report results.bin

echo "GET http://www.vol-4-ok.ru:8080/api/adverts/?userId=2&deleted=1" | vegeta attack -name=50qps -rate=50 -duration=3s --output results.bin;

echo "GET http://www.vol-4-ok.ru:8080/api/adverts/?userId=2&deleted=1&count=20&startId=1" | vegeta attack -name=50qps -rate=10 -duration=15s --output results.bin;
echo "GET http://www.vol-4-ok.ru:8080/api/profile/3" | vegeta attack -name=50qps -rate=10 -duration=15s --output results.bin;

