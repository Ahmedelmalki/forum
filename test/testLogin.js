const http = require('http');

const options = {
    hostname: 'localhost',
    port: 9898,
    path: '/login/submit',  
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    }
};

const requestsPerBatch = 300;  // Number of simultaneous requests
const totalBatches = 50;        // Number of batches to send
const delayBetweenBatches = 1000; // Delay between batches in ms

function makeRequest() {
    return new Promise((resolve) => {
        const req = http.request(options, (res) => {
            resolve(res.statusCode);
        });

        req.on('error', (error) => {
            resolve(error.code);
        });

        // Sample login payload
        const data = JSON.stringify({
            email: 'test@test.com',
            password: 'password123'
        });

        req.write(data);
        req.end();
    });
}

async function runTest() {
    for (let batch = 0; batch < totalBatches; batch++) {
        console.log(`Sending batch ${batch + 1}/${totalBatches}`);

        const promises = Array(requestsPerBatch).fill()
            .map(() => makeRequest());

        const results = await Promise.all(promises);

        // Count status codes
        const counts = results.reduce((acc, code) => {
            acc[code] = (acc[code] || 0) + 1;
            return acc;
        }, {});

        console.log('Results:', counts);

        // Wait before next batch
        if (batch < totalBatches - 1) {
            await new Promise(resolve => setTimeout(resolve, delayBetweenBatches));
        }
    }
}

runTest();

// node test/testLogin.js
