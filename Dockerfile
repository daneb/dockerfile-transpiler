
FROM ruby:3.0.3-alpine

RUN apk add --update --virtual \
  nodejs \
  yarn \
  readline \
  file \
  imagemagick \
  git \
  tzdata

WORKDIR /app
COPY . /app/

ENV BUNDLE_PATH /gems
RUN yarn install
RUN gem install bundler:2.2.32
RUN bundle install

ENTRYPOINT ["bin/rails"]
CMD ["s", "-b", "0.0.0.0"]

EXPOSE 3000
