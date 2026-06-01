# Parse REDIS_URL into REDIS_PASSWORD (and REDIS_HOST, REDIS_PORT, REDIS_DB when needed).
# Usage: set REDIS_URL, then: . scripts/parse-redis-url.sh

_parse_redis_url() {
	if [ -z "${REDIS_URL:-}" ]; then
		echo "REDIS_URL is required" >&2
		return 1
	fi

	url="$REDIS_URL"
	case "$url" in
	redis://*) rest="${url#redis://}" ;;
	rediss://*) rest="${url#rediss://}" ;;
	*)
		echo "REDIS_URL must use redis:// or rediss://" >&2
		return 1
		;;
	esac

	# redis://[:password@]host:port/db
	case "$rest" in
	*@*)
		auth="${rest%%@*}"
		hostpart="${rest#*@}"
		;;
	*)
		auth=""
		hostpart="$rest"
		;;
	esac

	# auth is :password or empty
	case "$auth" in
	:*)
		REDIS_PASSWORD="${auth#:}"
		;;
	"")
		REDIS_PASSWORD=""
		;;
	*)
		REDIS_PASSWORD="$auth"
		;;
	esac

	hostpath="${hostpart%%/*}"
	REDIS_DB="${hostpart#*/}"
	[ "$REDIS_DB" = "$hostpart" ] && REDIS_DB="0"
	REDIS_DB="${REDIS_DB%%\?*}"

	case "$hostpath" in
	*:*)
		REDIS_HOST="${hostpath%%:*}"
		REDIS_PORT="${hostpath#*:}"
		;;
	*)
		REDIS_HOST="$hostpath"
		REDIS_PORT="6379"
		;;
	esac

	export REDIS_PASSWORD REDIS_HOST REDIS_PORT REDIS_DB
}

# Rewrite host/port; keeps password, DB index, and rediss:// if present.
redis_url_rewrite_host() {
	url="$1"
	new_host="$2"
	new_port="$3"

	scheme="redis"
	case "$url" in
	rediss://*)
		scheme="rediss"
		rest="${url#rediss://}"
		;;
	redis://*)
		rest="${url#redis://}"
		;;
	*)
		echo "invalid REDIS_URL" >&2
		return 1
		;;
	esac

	_old="$REDIS_URL"
	REDIS_URL="$url"
	_parse_redis_url || return 1
	REDIS_URL="$_old"

	auth=""
	if [ -n "$REDIS_PASSWORD" ]; then
		auth=":${REDIS_PASSWORD}"
	fi

	printf '%s://%s@%s:%s/%s' "$scheme" "$auth" "$new_host" "$new_port" "$REDIS_DB"
}

_parse_redis_url
