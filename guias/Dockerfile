FROM hhvm/hhvm-proxygen:latest

# Remove default files
RUN rm -rf /var/www
RUN mkdir /var/www

WORKDIR /var/www

# Add project
ADD hh_autoload.json .

# Run composer
ADD composer.phar .

RUN php composer.phar require hhvm/hsl hhvm/hhvm-autoload

# Add our files
ADD . .

# Run type checker
RUN hh_client; exit 0

# Run autoload finder
RUN hhvm vendor/bin/hh-autoload.hack

EXPOSE 80