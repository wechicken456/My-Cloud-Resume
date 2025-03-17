# Cloud Resume

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
aws s3 cp ./dist s3://tin-cloud-resume/ --recursive
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