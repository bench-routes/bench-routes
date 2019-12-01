const webdriver = require('selenium-webdriver'),
    chrome = require('selenium-webdriver/chrome'),
    By = webdriver.By,
    Key = webdriver.Key

{
    describe, it
} require('selenium-webdriver/testing')

require('chromedriver')

let chromeOptions = new chrome.Options(),
    driver;

chromeOptions.addArguments('--no-sandbox');
chromeOptions.addArguments('--disable-dev-shm-usage');
chromeOptions.addArguments('--disable-gpu');

describe('launch services via virtual browser test', function() {
    this.timeout(200000);
    it('launching chrome browser', (done) => {
        driver = new webdriver
            .Builder()
            .setChromeOptions(chromeOptions)
            .forBrowser('chrome')
            .build();
        driver.then(() => {
            done();
        });
    });
    it('load test file', (done) => {
        driver.get('http://localhost:9090/test').then(() => {
            done();
        });
    });
    it('check if test page loaded', (done) => {
        driver.findElement(By.xpath('/html/body/button[1]')).then(() => {
            driver.findElement(By.xpath('/html/body/button[5]')).then(() => {
                done();
            }).catch(e => {
                throw e;
            });
        }).catch(e => {
            throw e;
        });
    });
    describe('trigger services', function() {
        this.timeout(50000);
        it('ping', done => {
            driver.findElement(By.xpath('/html/body/button[1]')).click().then(() => {
                done();
            }).catch(e => {
                throw e;
            });
        });
        it('flood-ping', done => {
            driver.findElement(By.xpath('/html/body/button[3]')).click().then(() => {
                done();
            }).catch(e => {
                throw e;
            });
        });
        it('jitter', done => {
            driver.findElement(By.xpath('/html/body/button[5]')).click().then(() => {
                done();
            }).catch(e => {
                throw e;
            });
        });
        it('req-res-delay', done => {
            driver.findElement(By.xpath('/html/body/button[7]')).click().then(() => {
                done();
            }).catch(e => {
                throw e;
            });
        });

        it('get-route-details', done => {
            driver.findElement(By.xpath('/html/body/button[9]')).click().then(() => {
                done();
            }).catch(e => {
                throw e;
            });
        });
        it('flood-ping', done => {
            driver.findElement(By.xpath('/html/body/button[10]')).click().then(() => {
                done();
            }).catch(e => {
                throw e;
            });
        });
        it('jitter', done => {
            driver.findElement(By.xpath('/html/body/button[11]')).click().then(() => {
                done();
            }).catch(e => {
                throw e;
            });
        });
        it('req-res-delay', done => {
            driver.findElement(By.xpath('/html/body/button[12]')).click().then(() => {
                done();
            }).catch(e => {
                throw e;
            });
        });
        it('jitter', done => {
            driver.findElement(By.xpath('/html/body/button[13]')).click().then(() => {
                setTimeout(() => {
                    done();
                }, 10000);
            }).catch(e => {
                throw e;
            });
        });
    });
    describe('shutting down services', function() {
        this.timeout(50000);
        it('ping', done => {
            driver.findElement(By.xpath('/html/body/button[2]')).click().then(() => {
                done();
            }).catch(e => {
                throw e;
            });
        });
        it('flood-ping', done => {
            driver.findElement(By.xpath('/html/body/button[4]')).click().then(() => {
                done();
            }).catch(e => {
                throw e;
            });
        });
        it('jitter', done => {
            driver.findElement(By.xpath('/html/body/button[6]')).click().then(() => {
                done();
            }).catch(e => {
                throw e;
            });
        });
        it('req-res-delay', done => {
            driver.findElement(By.xpath('/html/body/button[8]')).click().then(() => {
                done();
            }).catch(e => {
                throw e;
            });
        });
    });
    describe('closing the browser', function() {
        this.timeout(3000);
        it('closing browser service', (done) => {
            driver.quit();
            done();
        });      
    });
});
