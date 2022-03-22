# 💬 Epigram 

> **[ep·i·gram][wikipedia]** *noun*
: a pithy saying or remark expressing an idea in a clever and amusing way.

Epigram is a simple web service for communities to immortalize the enlightening, funny, or downright dumb quotes that they hear.


[wikipedia]: https://en.wikipedia.org/wiki/Epigram

## Features

- [x] Users can submit and view quotes.
- [x] Quotes are organized in chronological order, and in sections by year.
- [x] Authorization is delegated to a configurable OpenID Connect provider.
- [x] Access restricted to only those who correctly answer a few questions.
- [ ] Dark mode support.
- [ ] Expanded admin control functions.

## Project Status

Epigram still under development, and should be considered a work in progress.

## Usage

### Installation

Epigram can be compiled and installed directly from source as follows:

```bash
go install github.com/willbicks/epigram@latest
```

Alternatively, Docker container images are available at [ghcr.io/willbicks/epigram](https://ghcr.io/willbicks/epigram).

### Configuration

While most configuration parameters have a default value, 

Epigram expects a `config.yml` in the same directory as it is run, or at the location specified by the `EP_CONFIG` environment variable.

TODO: Explain config file contents. For more information, see [Configuration](docs/config.md).

## Documentation

- [Project Structure / Architecture](docs/structure.md)
- [Configuration](docs/config.md)

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss your proposed change.

Please make sure to update tests as appropriate.

## Copyright and License

Copyright (c) 2022 [Will Bicks](https://www.willbicks.com)

Distributed under the BSD 3-Clause License. For more details, see [LICENSE](LICENSE).
