# Exitstatus

**Concepts**

db
  * db is file-based, JSON-based
  * checks get updated in the file as they check in

**Main threads**

manager
  * manages main go channels 
  * manages the logger thread
  *
  * gracefully handles signals

**worker**
  * checks last checkin key
  * updates the status
  * publishes message to interested parties

**API server**
  * starts the http-based api server

**alerter**

**pub/sub**
check producer - the workers
check consumers - SMS, email, ledisDB updating

**config file**

* jobs
  * last time they checked in
  * last time they were checked
  * interval for checking
  * enabled / disabled
  * status
    * healthy - ping heard within time window
    * unhealthy - did not hear ping within time window
    * unknown - new ping set up, or ping is disabled
* email config
  * smtp credentials
* sms config

## executable

* `exitstatus --serve` - run the client
* `exitstatus` - run the client

    checkup add check 60 s3_backup "server backup to S3"
      please ping this URL with curl: https://checkup.sh/check/aasdfkjdkfj

    checkup status s3_backup
      enabled: true
      status: HEALTHY
      last heard from: 13 minutes ago
      interval: 15 minutes
      next check: 2 minutes from now


alerters:
  - name: blah
    type: webhook
    url: https://blah.com/ping
  - name: email ping
    type: mail
    smtp_creds: 

# check states

* healthy - a service has checked in before the check window expiry
* unhealthy - a service has failed to check in 
* unknown - exitstatus was started after the check window expiry

# libraries

* https://github.com/gorhill/cronexpr
