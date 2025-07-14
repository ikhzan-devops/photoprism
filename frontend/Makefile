# Copyright © 2018 - 2025 PhotoPrism UG. All rights reserved.
#
# Questions? Email us at hello@photoprism.app or visit our website to learn
# more about our team, products and services: https://www.photoprism.app/

# Declare "make" targets.
all: install build
dep: install-npm install-testcafe install
dep-list:
	npx npm-check-updates
notice:
	@echo "Creating license report for frontend dependencies..."
	license-report --only=prod --config=.report.json > NOTICE
install-npm:
	sudo npm install --unsafe-perm=true --allow-root -g npm@latest npm-check-updates@latest license-report@latest
install-testcafe:
	npm install -g testcafe@latest
install-eslint:
	npm install -g eslint globals @eslint/eslintrc @eslint/js eslint-config-prettier eslint-formatter-pretty eslint-plugin-html eslint-plugin-import eslint-plugin-node eslint-plugin-prettier eslint-plugin-promise eslint-plugin-vue eslint-webpack-plugin vue-eslint-parser prettier
npm-install:
	$(info Installing dependencies...)
	npm install --no-update-notifier --no-audit
install: npm-install
npm-update:
	$(info Updating dependencies...)
	npm update --save --package-lock
update: npm-update npm-install
upgrade: update
watch:
	npm run watch
build:
	npm run build
lint:
	npm run lint
fmt:
	npm run fmt
test:
	npm run test
testcafe:
	npm run testcafe
acceptance-local:
	npm run acceptance-local
gettext-extract:
	npm run gettext-extract
gettext-compile:
	npm run gettext-compile

# Declare all targets as "PHONY", see https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html.
MAKEFLAGS += --always-make
