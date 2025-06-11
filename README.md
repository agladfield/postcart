<!-- TITLE: Postc.art: Send Beautiful Postcards With Email! -->

This is a submission for the [Postmark Challenge: Inbox Innovators](https://dev.to/challenges/postmark).

## What I Built

Postc.art enables you to send beautiful postcards straight from email. No need to remember any passwords or sign into a website. Just fill out an email with the standardized template and send!


![Demo Loop](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/16ufzrh41egv6ee9kpin.gif)

Your recipient have terrible internet connection, an outdated computer, or doesn't trust HTML email content? No problem! Each postcard email also comes with a text/ASCII version of the email to ensure your recipient at least has a vague idea of what your postcard looks like and what you wanted to tell them.


![ASCII Loop](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/l9bs9b6huyuy84rws5ns.gif)

(ASCII Version of your postcards!)

## Demo

To try Postc.art out you have two choices:

### Live Service

You can try out sending a postcard right now straight from your email!

Create a new email, address it to `send@postc.art`, give it a subject, and then implement this template:

```
To: Name <email@domain.com>
From: Name <email@domain.com>
Artwork: { City | Islands | Lake | Mountains }
Style: { Illustrated | Painting | Photo | Vintage Photo }
Border: { None | Stripes | Lines | Cubes | Artwork }
Font: { Marker | Polite | Retro | Typewriter }
Textured: { Yes | No }
Country: { The ISO2, ISO3, or most common spelling of any recognized countries }
{ Message Content }
```

#### Use the Template

[Click here to draft an email with this template!](mailto:send@postc.art?subject=New%20Postcard&body=To:%20Name%20%3Cemail@domain.com%3E%0AFrom:%20Name%20%3Cemail@domain.com%3E%0AArtwork:%0AStyle:%0ABorder:%0AFont:%0ATextured:%0ACountry:%0A%0A)

- The only required fields is `To:`. The message starts when you are done with the fields. If you do not specify `From:` the sevice will automatically use your email and name from your email.
- The values within the `{}` curly braces are the presently available options.
- If you choose to send emails via the live service please note that you have a balance of **3** postcards before you are cut off as to avoid racking up my AI bill. Additionally, if you send a postcard and it bounces or is marked as spam you will no longer be able to send postcards.
- If you are still unsure look at some of the example emails provided in the first gif.

**Note:** If you clone the repository and run the code yourself you can use attachments as artwork so long as they are one of the following formats:

```
image/png
image/webp
image/jpeg
```

However, for safety, security, and legal reasons I did not enable attachment artwork in the live service.

### Running the code on your own

To run the code you MUST HAVE:

- A recent-ish version of Go installed
    > [Install Go Guide](https://go.dev/doc/install)
- LibVIPS (an image processing library) installed
    > [Install LibVIPS](https://www.libvips.org/install.html)

To run the project, clone or download the repository, cd into the src directory (`cd src`), run `go mod download` then you have two choices.

Ensure you are in the `src` directory before continuing otherwise Go will not be able to detect the project.

If you would just like to see the Postcard generation in action, you can run (AI features are disabled by default):

```
go run ./cmd/postcard-demo
```

Otherwise, if you prefer to run the entire server (webhook processing, image uploading, and email sending) from start to finish run:

```
go run ./cmd/postcart
```

To get the full server to run properly you need to pass at least these environment variables:

- `POSTMARK_SERVER_TOKEN`: The token for the Postmark server you want to use
- `POSTMARK_AUTH_USER`: The webhook basic authentication username to protect your server
- `POSTMARK_AUTH_PASS`: The webhook basic authentication password to protect your server
- `POSTMARK_INBOUND_EMAIL`: The inbound email for which to expect postcards to be addressed to
- `POSTMARK_EMAIL_DOMAIN`: The domain from which to send out other emails on

Optional Env Values Include:

- `GCP_CRED_PATH`: The path to your Google Cloud/Vertex AI credentials needed for using Google's Imagen4 image generation AI. See the guide on how to get it.

    > [Creating Service Account for Vertex AI Guide](https://docs.mindmac.app/how-to.../add-api-key/create-google-cloud-vertex-ai-api-key)

    > IMPORTANT: In addition to giving the service account `Vertex AI Service Agent` it must also have a Storage administrative role of some kind that it may to write to buckets with your project if you wish to use the longer storage.

- `GCP_PROJECT`: The project you created your service account above with. Make sure it has Vertex AI enabled.
- `GCP_BUCKET`: A Google Cloud bucket path for persisting image uploads. If not included `tmpfiles.org` will be used instead which only stores the images for 60 minutes.
- `PORT`: Overrides the http server port (defaults to 8080)

Env Toggles:

- `INSTALL_FONTS`: If set to `true`, will install fonts into the current user (you's) Font directory (untested on Windows but should work). Otherwise, your computer will rollback to seriff fonts if these are not found.
- `USE_AI`: If set to `true`, will use Google's Imagen AI to generate artwork for postcards. Must be paired with `VERTEX_AI_KEY`.
- `ALLOW_ATTACHMENTS`: If set to `true`, will enable the processing of attachment images as postcard artwork.

While GoLang does offer environment variable libraries I did not feel those necessary for this project, as such you can find all of the environment variables and logic inside `src/pkg/shared/env/env.go`.

Here is a template you can copy and paste to put your environment variables in:

```
export POSTMARK_SERVER_TOKEN={SERVER_TOKEN}
export POSTMARK_AUTH_USER={SERVER_BASIC_AUTH_USERNAME}
export POSTMARK_AUTH_PASS={SERVER_BASIC_AUTH_PASSWORD}
export POSTMARK_INBOUND_EMAIL={INBOUND_EMAIL}
export POSTMARK_EMAIL_DOMAIN={EMAIL_DOMAIN}
export GCP_CRED_PATH=""
export GCP_PROJECT=""
export GCP_BUCKET=""
USE_AI=false
INSTALL_FONTS=false
ALLOW_ATTACHMENTS=true
```

You will know that everything configured and started correctly if you end up with this screen:

![Good to go splash](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/ywpyt4izpwkd7oacxw61.gif)

If configured correctly, your Postmark server should have a new standard template in it for sending postcards called `Deliveries Template`.

## Code Repository

[Github Repository](https://github.com/agladfield/postcart) (Should Completely work so long as you followed the Go and LibVIPS installation processes)

## How I Built It

### Implementation Process

When I first received the prompt for the challenge I took a little while to think about what sort of information could be conveniently shared via email in a structured format. That is when I landed on the idea that all of the information contained on a postcard could easily be fitted to a standardized email.

I decided that, since either end-user (the sender or the recipient) deal with the postcards via email only there was no neeed for a traditional front end/website. So first I came up with the design and resources of the postcards (found in the `src/pkg/cards/res` folder). Then I set about writing the code to assemble the completed postcard images (& later text/ASCII). Once I was satisfied with how they look and were assembled, I moved on to integrating with Postmark.

As there is no officially maintained Postmark API Library for GoLang I had to make my own. So included in my code is one I made my own that only includes the features and types used in this program. For instance, I did not implement the default send email endpoint as I only send emails/postcards using templates. Something important to note for anyone unfamiliar with Go evaluating the Postmark code I provided should know is that in Go the names/keys of a struct's properties do not necessarily match up 1 for 1 like they would in say Typescript. Instead, Go uses what are called JSON tags to identify under what string key it should look for a value when turning raw bytes into a structured type. While conveniently Postmark's JSON keys basically universally match up with what Go expects in a struct declaration (PascalCase), there were a few small semantics where the JSON tag differs from the field key in my code such for example: HTMLBody string `json:"HtmlBody"`. Go's styling guidelines/LSPs strongly emphasize keeping acronyms capitalized so I followed that convention throughout the program. Having gotten that out of the way, I implemented types for the webhook events I want to receive and process (inbound emails, delivered, bounce, and spam complaint notifications). I then implemented the request and response types and http calls for the API routes I intended to use (templates, inbound rules, send with template).

Now that I had interfacing with Postmark established, I built the webhook server and inbound email parser. Once I got the parser working and placing the information on the postcards I setup uploading the images to third party services. This way I can provide URLs to display them in emails as I found that they would get cut off if I tried to embed them. To make them work in a cohesive fashion I use Postmark's template feature which allows me to send less information with my email requests and have the postcard images consistently fit edge to edge. All-in-all I am happy with how it turned out.

### Tech Stack

- **Go Lang:** For robust error handling, performance, large standard library (http server & clients included), and static compilability.
- **LibVIPS:** A high performance image processing library written in C.
- **Postmark:** for email processing, status notifications, templating, and sending.
- **Google's Imagen4:** for generating beautiful artwork images for the Postcards.
- **Google's Cloud Storage:** for permanent image storage and embedding.
- **tmpfiles.org:** for quick and easy image upload for embedding the images into sent emails.

### Postmark experience

This is my first full implementation of Postmark, I have only previously used it briefly in testing where I observed solid deliverability. Having gotten to test and implement a greater breadth of their services while building this I am (rather unsurprisingly) increasingly pleased with Postmark's ease of use and reliability of their documentation. The only time I ran into any issues with Postmark while building my project was a result of my incorrect implementation and not anything on their end. I am grateful to have chosen Postmark as an email provider because otherwise I would not have heard about this competition and would not have had this idea!

### Ways I would improve upon my submission

- **_Adding an actual database:_** I had originally developed this with Turso/LibSQL but as development progressed the network latency became unbearable and so I had to scrap using it. Moving forward I would probably pair it with something like PostgreSQL.
- **_Splitting up the Postcard Images:_** It is quick and easy joining the postcard images completely but it would look better if they were rotated and partially overlapping one another for a more realistic, authentic look.
- **_Increasing Option Variety:_** Expanding configuration offerings. While I am satisfied with the four style, artwork, and font options, I would definitely add many more moving forward.
- **_Adding randomization:_** Adding random imperfections to the cards to make them feel less cookie cutter.
- **_Instrument it:_** Right now there's no observability metrics provided by the application. The only way to see how things are going is via the Postmark dashboard. I would implement Prometheus metrics to get a better understanding of performance and usage patterns.
- **_URL Permanance:_** The maximum time a signed URL may expire with Google Cloud Storage is 7 days from issuing. This is good for most recipients but in case someone wanted to check on an old postcard they received the embedded URL should last longer than a week.
- **_Smarter Parser:_** While my current email parser is reasonably fast and reliable, it is rather rigid and I think if I were to do it over I would completely redo the message/text body parser.

### Closing note

This is my first ever public release of a programming project I have written. Paid or otherwise. I would appreciate criticisms and other feedbacks as I continue to try to improve. Thanks for your time!

Additionally, I love this idea so that much that I will be turning it into an actual business. So stay tuned for here or at the website ([https://postc.art](https://postc.art)) to see it in action.

### Edit 6/10:

For some reason the example gifs that should be loaded from my site were not loading and so I have fixed this by uploading them here to dev.to.
