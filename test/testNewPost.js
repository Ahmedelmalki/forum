const http = require('http');
const querystring = require('querystring');

const options = {
    hostname: 'localhost',
    port: 9898,
    path: '/newPost',
    method: 'POST',
    headers: {
        'Content-Type': 'application/x-www-form-urlencoded', // Form data
        'Cookie': 'forum_session=3d2bdc23-9156-42bb-829b-67bf801d5f07', // Simulate an authenticated session
    }
};

// Test parameters
const requestsPerBatch = 10;  // Number of simultaneous requests
const totalBatches = 10;       // Number of batches to send
const delayBetweenBatches = 1000; // Delay between batches in ms

function makeRequest() {
    return new Promise((resolve) => {
        const req = http.request(options, (res) => {
            resolve(res.statusCode);
        });

        req.on('error', (error) => {
            resolve(error.code);
        });

        // Sample form data payload
        const data = querystring.stringify({
            title: 'Test Title',
            content: 'This is a test post content.',
            'categories[]': 'tech' // Form array (use correct keys for your handler)
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
