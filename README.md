# ten34

A globally-distributed, eventually-consistent, 100% available key-value store.

> Route 53 isn't really a database, but then again, neither is Redis.

_[Corey Quinn](https://twitter.com/QuinnyPig/status/1173371936342044672)_

## Usage

### Creating a database

```shell
ten34 createdb route53://my.db
```

### Deleting a database

```shell
ten34 dropdb route53://my.db
```

### Setting a key

```shell
ten34 -d route53://my.db put foo bar
```

### Getting a key

```shell
ten34 -d route53://my.db get foo
```

### Deleting a key

```shell
ten34 -d route53://my.db del foo -d
```

## Development

    $ make build

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/craftyphotons/ten34. This project is intended to be a safe, welcoming space for collaboration, and contributors are expected to adhere to the [Contributor Covenant](http://contributor-covenant.org) code of conduct.

## License

The gem is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).

## Code of Conduct

Everyone interacting in the ten34 project’s codebases, issue trackers, chat rooms and mailing lists is expected to follow the [code of conduct](https://github.com/craftyphotons/ten34/blob/master/CODE_OF_CONDUCT.md).

## Additional Disclaimer

In addition to the terms of the MIT license, this project and its maintainers shall not be held responsible for costs and repercussions resulting from its use. This includes but is not limited to account closure by your cloud service provider for violating their terms of service and the disappointment of your peers for usage of this project for your actual database.