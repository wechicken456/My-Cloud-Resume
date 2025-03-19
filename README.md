# Cloud Resume

You can access my domain and resume [here](https://www.pwnph0fun.com).


## Project structure

```
./
|-- public/
│   |- index.html
|-- src/
│   |- data/
│   │   |- resumeData.ts
│   |- main.ts
│   |- style.css
|-- package.json
|-- tsconfig.json
|-- vite.config.ts
|-- tailwind.config.js
```

## Step 1 - HTML, CSS, Typescript

Standard stuffs.

Create app:
```bash
npm create vite@latest cloud-resume
npm install tailwindcss @tailwindcss/vite
npm install typescript
```

Then test the project with
```bash
npm run dev
```

## Step 2 - Create S3 bucket.

Some important configs to look for when creating the bucket:


- Uncheck Block all public access to allow public access (required for static hosting).
- Enable static website hosting from the `Properties` tab of the bucket.

Then upload all the files in the production folder `dist/` to the bucket from the terminal:
```bash
aws s3 cp ./dist s3://pwnph0fun-cloud-resume/ --recursive
```

Then every time we change the front end, 
Remember to **Invalidate** the CDN cache:
```
aws cloudfront create-invalidation --distribution-id <your-cloudfront-id> --paths "/*"
```

## Step 3 - Create Origin Access Identity

[Here](https://docs.aws.amazon.com/AmazonS3/latest/userguide/tutorial-s3-cloudfront-route53-video-streaming.html#cf-s3-step3)

By using an OAI, you make sure that viewers can't bypass CloudFront and get the video directly from the S3 bucket. Only the CloudFront OAI can access the file in the S3 bucket. 

1. In the *left navigation plane* of AWS CloudFront console, choose **Origin Access**.
2. Under the **Identities** tab, select **Create origin access identity**.
3. Enter a name for it. 
4. Create

## Step 4 - Create CloudFront Distribution

[Here](https://docs.aws.amazon.com/AmazonS3/latest/userguide/tutorial-s3-cloudfront-route53-video-streaming.html#cf-s3-step4)

Some important configs to look for when creating the distribution:

- Origin domain: Select your S3 bucket from the dropdown (e.g., tin-vuong-resume-<random-string>.s3.amazonaws.com). Do **NOT** choose the *website endpoint* option, even if AWS recommends it.
- Leave Origin path blank.
- For **Origin Access**, choose the **Origin access control settings (recommended)**.
    - Select the created OAI from step 3.
    - Copy policy and paste it into the policy of the S3 bucket.
- Redirect HTTP to HTTPS.
- Default cache behavior settings: Leave as default (GET, HEAD allowed).
- Default root object: Enter `index.html`.
- Price class: only North America and Europe.
- Leave Web Application Firewall (WAF) as Do not enable (unless you need it).

Remember to **Invalidate** the CDN cache every time we change the **front end**:
```
aws cloudfront create-invalidation --distribution-id <your-cloudfront-id> --paths "/*"
```

## Step 5 - Register a custom domain with Route53

[Here](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/domain-register.html#domain-register-procedure-section)

I registered a domain called `pwnph0fun.com`.
***IMPORTANT***: For Route 53 public and private DNS and health checks, the control plane is located in the us-east-1 AWS Region and the data planes are globally distributed.


### Substep 1 - Request a SSL certificate for your viewers to use HTTPS

[Here](https://docs.aws.amazon.com/AmazonS3/latest/userguide/tutorial-s3-cloudfront-route53-video-streaming.html#cf-s3-step6-create-SSL)

***IMPORTANT***: make sure your region is set to `us-east-1` because AWS CF is a **GLOBAL** region that uses `us-east-1` as its default and the ACM SSL certificate has to be created from the same region.

1. Go to AWS Certificate Manager (ACM).
2. Request a public certificate.
3. Enter `*.pwnph0fun.com` in the Add Domain Names section.
4. Select DNS Validation.
5. Create the request.
6. The request status will be `Pending`, go ahead and click on `Create records in Route53`.
7. The request status will turn to `Success` soon.

Why do we need to add a CNAME record? Didn't we just request the certificate? 

- When you request an SSL/TLS certificate from ACM, AWS needs to confirm that you have the authority to use the domain name you're requesting the certificate for.
- DNS validation involves adding a specific CNAME (Canonical Name) record to your domain's DNS settings.
- ACM then checks for the presence of this CNAME record to verify your domain ownership BEFORE validating your SSL certificate request for your domain.

The **unique** CNAME record acts as proof that you have **ownnership** over your domain's DNS settings.
Because only someone with access to your DNS settings can add this record, its presence verifies your authority.

>In Simple Terms: Imagine ACM sending you a secret code (the CNAME record). To prove you own your house (your domain), you need to write that code on your front door (your DNS settings). ACM then checks the front door to see if the code is there. If it is, they know it's your house.


### Substep 2 - Add your custom domain to CF

[Here](https://docs.aws.amazon.com/AmazonS3/latest/userguide/tutorial-s3-cloudfront-route53-video-streaming.html#cf-s3-step6-create-SSL)

1. Go to **CloudFront** management, edit **Settings** to add **Alternate domain name (CNAME)**.
2. I chose mine to be ```www.pwnph0fun.com```.
3. For **Custom SSL certificate**, choose the one we created above.
4. Save changes, wait for the **Last modified** status to change from **Deploying** to a timestamp.

### Substep 3 - Create a DNS record to route traffic from your alternate domain name to your CloudFront distribution's domain name

[Here](https://docs.aws.amazon.com/AmazonS3/latest/userguide/tutorial-s3-cloudfront-route53-video-streaming.html#cf-s3-step6-DNS-record)

1. Go to **Route53** management -> **Hosted Zones** -> click on the created domain.
2. Add the *subdomain*. e.g. we used ```www.pwnph0fun.com``` in the substep above, so I'll add `www` here. 
3. Enable **Alias**.
4. Route traffic to **Alias to CloudFront distribution**.
5. Choose the CF that we created.
6. Create record.


## Step 6 - Create a DynamoDB

[Here](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/getting-started-step-1.html)

[DynamoDB structure](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.CoreComponents.html?icmpid=docs_dynamodb_help_panel_hp_table#HowItWorks.CoreComponents.TablesItemsAttributes)

We can NOT simply use "Count" as the primary key because:

If `Count` is the primary key, you need to know the exact value of `Count` to fetch it (e.g., `Count` = 0). However, for a visitor counter, you want to increment a count, not fetch a static key. Using `Count` as the primary key doesn’t make sense for this use case because:
You can’t easily update a primary key (you’d need to delete and re-insert the item).
You likely want a single item with a fixed key and an attribute that increments (e.g., `Count` as an attribute, not the key).

**Solution**: Redesign the table to have a fixed primary key (e.g., `ID` as a String) and store the counter as an attribute (e.g., `Count` as a Number). Then use `GetItem` or `UpdateItem` to fetch/increment it.


Something like this:

```bash
aws dynamodb create-table --table-name VisitorCounter \ 
    --attribute-definitions AttributeName=ID,AttributeType=S \ 
    --key-schema AttributeName=ID,KeyType=HASH \ 
    --billing-mode PAY_PER_REQUEST \ 
    --table-class STANDARD
```

Then:
```bash
aws dynamodb put-item --table-name VisitorCounter \
    --item '{"ID": {"S": "visitor"}, "Count": {"N": "0"}}'
```

Then code up the GetCount() and IncrementCount() using the AWS SDK for Go v2.


## Step 7 - Set up Lambda

We will setup the backend for the visitor counter first, then test it with `cURL`, then add frontend Typescript to communicate with the backend.

[Tutorial](https://docs.aws.amazon.com/lambda/latest/dg/services-apigateway-tutorial.html)

### Substep 1 - Create a policy for Lambda to interact with DynamoDB.

### Substep 2 - Create an execution role for Lambda to interact with DynamoDB.

Use the same policy created in Substep 1.

Choose **Lambda** for *use case*. Add the policy created in Substep 1.

Then to integrate it with AWS API Gateway, see [CLI version](https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-proxy-integration-using-cli.html) and [console version](https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-create-api-as-simple-proxy-for-lambda.html).

Proxy integration: we handle the HTTP request with headers, and we have to specify the response format as well.
Custom integration: The headers are abstracted, we only need to specify the mapping from the HTTP request's input data into what our Lambda wants to see. But this involves more setup. In this project, we used Proxy Integration.


### Substep 3 - Create the Lambda function

```
GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap main.go
zip myFunction.zip bootstrap
```

## Step 8 - Setup API Gateway

### Substep 1 - Create the `/api` resource on API Gateway
The `/api` prefix groups your endpoints, which is useful if you plan to add more API endpoints later (e.g., `/api/stats`, `/api/resetCount`). It keeps your API distinct from other potential resources (e.g., `/web`, `/admin`).

Convention: Many REST APIs use `/api` as a standard prefix to indicate an API endpoint, improving readability and aligning with common practices.

### Substep 2 - Create the `/getCount` and `/incrementCount` resource on API Gateway

- Choose the Lambda Integration option.
- Enable the CORS option.
- Enable Proxy Integration.
- Choose our Lambda function.

The `OPTIONS` resoruce lets browswer find out about CORS permissions: 
Browsers send an `OPTIONS` preflight request for `POST /putCount` (and potentially `GET /getCount` if headers change) to check CORS permissions. Your Lambda sets CORS headers (`Access-Control-Allow-Origin`: *), but API Gateway needs `OPTIONS` methods to respond to preflights.

I edited the methods allowed in CORS to only `GET, POST, OPTIONS` since we only need `GET` for `getCount` and `POST` for `incrementCount`


## Step 9 - Create Front-End code to fetch the APIs

See `main.ts`

Fetch once on DOM load, then fetch again every 5 seconds.

Use `localStorage` to maintain per browser increment: don't increment on reloads within the same browser.




