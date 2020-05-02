const webdriver = require('selenium-webdriver'),
    chrome = require('selenium-webdriver/chrome'),
    By = webdriver.By,
    Key = webdriver.Key

require('chromedriver')

let chromeOptions = new chrome.Options(),
    driver;

chromeOptions.addArguments('--no-sandbox');
chromeOptions.addArguments('--disable-dev-shm-usage');
chromeOptions.addArguments('--disable-gpu');

describe('launch Services via virtual browser test', function() {
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
        // Starts ping and jitter using startPing and startJitter
        driver.findElement(By.xpath('/html/body/div/div[1]/div[1]/div/button[1]')).then(() => {
            driver.findElement(By.xpath('/html/body/div/div[1]/div[3]/div/button[1]')).then(() => {
                done();
            }).catch(e => {
                throw e;
            });
        }).catch(e => {
            throw e;
        });
    });
    describe('trigger Services', function() {
        this.timeout(50000);
        it('ping', done => {
            // starts Services again
            driver.findElement(By.xpath('//*[@id="btnGroupDrop1"]')).click().then(() => {
                driver.findElement(By.xpath('/html/body/div/div[1]/div[1]/div/button[1]')).click().then(() => {
                    done();
                }).catch(e => {
                    throw e;
                });
            }).catch(e => {
                throw e;
            });
        });
        it('flood-ping', done => {
            driver.findElement(By.xpath('//*[@id="btnGroupDrop2"]')).click().then(() => {
                driver.findElement(By.xpath('/html/body/div/div[1]/div[2]/div/button[1]')).click().then(() => {
                    done();
                }).catch(e => {
                    throw e;
                });
            }).catch(e => {
                throw e;
            });
        });
        it('jitter', done => {
            driver.findElement(By.xpath('//*[@id="btnGroupDrop3"]')).click().then(() => {
                driver.findElement(By.xpath('/html/body/div/div[1]/div[3]/div/button[1]')).click().then(() => {
                    done();
                }).catch(e => {
                    throw e;
                });
            }).catch(e => {
                throw e;
            });
        });
        it('req-res-delay', done => {
            driver.findElement(By.xpath('//*[@id="btnGroupDrop4"]')).click().then(() => {
                driver.findElement(By.xpath('/html/body/div/div[1]/div[4]/div/button[1]')).click().then(() => {
                    done();
                }).catch(e => {
                    throw e;
                });
            }).catch(e => {
                throw e;
            });
        });

        // Get route details
        it('get-route-details', done => {
            driver.findElement(By.xpath('//*[@id="btnGroupDrop5"]')).click().then(() => {
                driver.findElement(By.xpath('/html/body/div/div[2]/div/div/button')).click().then(() => {
                    done();
                }).catch(e => {
                    throw e;
                });
            }).catch(e => {
                throw e;
            });
        });
        it('flood-ping', done => {
            driver.findElement(By.xpath('//*[@id="btnGroupDrop6"]')).click().then(() => {
                driver.findElement(By.xpath('/html/body/div/div[3]/div[1]/div/button')).click().then(() => {
                    done();
                }).catch(e => {
                    throw e;
                });
            }).catch(e => {
                throw e;
            });
        });
        it('jitter', done => {
            driver.findElement(By.xpath('//*[@id="btnGroupDrop7"]')).click().then(() => {
                driver.findElement(By.xpath('/html/body/div/div[3]/div[2]/div/button')).click().then(() => {
                    done();
                }).catch(e => {
                    throw e;
                });
            }).catch(e => {
                throw e;
            });
        });
        it('req-res-delay', done => {
            driver.findElement(By.xpath('//*[@id="btnGroupDrop8"]')).click().then(() => {
                driver.findElement(By.xpath('/html/body/div/div[3]/div[3]/div/button')).click().then(() => {
                    done();
                }).catch(e => {
                    throw e;
                });
            }).catch(e => {
                throw e;
            });
        });
        it('jitter', done => {
            driver.findElement(By.xpath('//*[@id="btnGroupDrop9"]')).click().then(() => {
                driver.findElement(By.xpath('/html/body/div/div[3]/div[4]/div/button')).click().then(() => {
                    setTimeout(() => {
                        done();
                    }, 10000);
                }).catch(e => {
                    throw e;
                });
            }).catch(e => {
                throw e;
            });
        });
    });
    describe('shutting down Services', function() {
        this.timeout(50000);
        it('ping', done => {
            driver.findElement(By.xpath('//*[@id="btnGroupDrop1"]')).click().then(() => {
                driver.findElement(By.xpath('/html/body/div/div[1]/div[1]/div/button[2]')).click().then(() => {
                    done();
                }).catch(e => {
                    throw e;
                });
            }).catch(e => {
                throw e;
            });
        });
        it('flood-ping', done => {
            driver.findElement(By.xpath('//*[@id="btnGroupDrop2"]')).click().then(() => {
                driver.findElement(By.xpath('/html/body/div/div[1]/div[2]/div/button[2]')).click().then(() => {
                    done();
                }).catch(e => {
                    throw e;
                });
            }).catch(e => {
                throw e;
            });
        });
        it('jitter', done => {
            driver.findElement(By.xpath('//*[@id="btnGroupDrop3"]')).click().then(() => {
                driver.findElement(By.xpath('/html/body/div/div[1]/div[3]/div/button[2]')).click().then(() => {
                    done();
                }).catch(e => {
                    throw e;
                });
            }).catch(e => {
                throw e;
            });
        });
        it('req-res-delay', done => {
            driver.findElement(By.xpath('//*[@id="btnGroupDrop4"]')).click().then(() => {
                driver.findElement(By.xpath('/html/body/div/div[1]/div[4]/div/button[2]')).click().then(() => {
                    done();
                }).catch(e => {
                    throw e;
                });
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
