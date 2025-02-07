# cool

## Important

This project is just a learning project, it is not meant to be used in production and is not provided somewhere (there is no public/private instance of this project).

There is no content from the providers in this repository.

A lot of considerations were made to maintain a low amount of network requests, so handle the main binary flags with care, and please, be mindful of the requests you make.

There is no better way of putting this, but the only supported provider is a... questionable website. What I am trying to say is that you shouldn't visit their homepage in public, you have been warned.

It was really funny fixing a bad css in a `<a>` while a video was playing muted in fullscreen.

## Why

I wanted to improve my knowledge in svelte and make a project that incorporates more than just a backend or a frontend, it ended up using a lot more than just that on the way, for example it had a htmx and a tauri frontend, it used sqlite instead of postgres, it had some powershell scripts.

It was also a good opportunity to handle large files, overall, I learned a lot.

## How

The main executable will fetch pages from the providers, parse the data, convert images to webp, create gif previews from videos, and save it to a postgres database.

The static website will fetch the data from the backend that serves not only the built static website but also content from the provider.

The API returns gzipped json and will always cache the response after successfully returning a response.

## Running

From the root directory, you can `go run .\cmd\backend` or `go run .\cmd\main`, for the frontend, you will need to `cd frontend` and either `npm run dev` or `npm run build`.

`main` is the project responsible for gathering the content of the providers, so you will need to run that first.

For dev, just `npm run dev` in the `frontend` directory and `go run .\cmd\backend` in the project root.

To run a "release" version of the project, `npm run build` the frontend then `go build/run` the backend.

`backend`:

- `-h2`: runs the server with HTTP/2, you will need a `cert.pem` and `key.pem` in the root directory.

`main`:

- `-provider`: The provider to scan.
- `-workers`: The number of concurrent workers to use.
- `-page`: The page to start at.
- `-per`: The number of posts per page.

## Interesting things

Using `cwebp -q 85` to convert all images reduced the size of the provider's content from 15gb to 10gb with no perceived quality loss.

The svelte frontend was SSR only, but I prefer static websites, so I ended up making a backend in Go that would serve the static website, its API and content of the providers.

There was logic to convert the image, check if the output (`webp`) was bigger than the input(`jpeg,png`), and if it was, use the input instead, this ended up being unnecessary since the savings were small, the cases where that happened weren't very common and when it did happen, it was only bigger for a few Kbs.

Dozens of videos autoplaying in the background can be expensive on the cpu! Never thought about that! I had a powershell script that would convert all videos to gifs by getting samples from the video and then concatenating them into a gif, it is now converted to a go function in the utils package.

Service workers are awesome, any request to the backend will be cached, and with the images also being cached, the website feels really snappy, not that it matters much since it is a very simple website.

The difference between http/1 and http/2 is insane, http/2 is really good for this project because of multiplexing.

Turns out, its very awkward to talk about this project to people.

## Ideas

- Maybe keeping tracking of the pages being fetched somewhere would be useful to resume the process and/or skip already fetched pages.
- A demo version of the website/database would be nice, not much time to do that.
- Some parts of the project would require a very small refactor to support multiple providers, just moving things to a map or helper functions to return db/data based on the provider.
- Keep track of the video.m3u8 stream-media-id -> filename in a file or database, this would be useful to avoid making a request for the m3u8.
- The download of media used to be handled in a goroutine, allowing for concurrent downloads, but I ended up removing it, maybe using a semaphore to limit the number of downloads(globally?) at a time would be useful.
- The way the website handles things while fetching the data could be improved, maybe add some spinners or something, the masonry feels like the slowest part of the website.
- There was a virtual list implementation I made, but I ended up removing it, maybe put it back so the only dependency could be removed.
- Maybe the logic of output being bigger than input should be used in the gif conversion, found one gif that was bigget than the video it was made from.
