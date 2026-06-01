# Resolve DATABASE_URL + DATABASE_USER/PASSWORD into RESOLVED_DATABASE_URL and DB_* parts.
# Usage: set vars, then: . scripts/parse-database-url.sh

_resolve_database_url_env() {
	if [ -z "${DATABASE_URL:-}" ]; then
		echo "DATABASE_URL is required" >&2
		return 1
	fi

	url="$DATABASE_URL"
	case "$url" in
	jdbc:*) url="${url#jdbc:}" ;;
	esac

	if [ -n "${DATABASE_USER:-}" ]; then
		case "$url" in
		postgresql://*|postgres://*)
			rest="${url#*://}"
			case "$rest" in
			*@*)
				# Full URL already (e.g. docker/compose.local.yml override)
				RESOLVED_DATABASE_URL="$url"
				;;
			*)
				query=""
				path="$rest"
				case "$rest" in
				*\?*)
					path="${rest%%\?*}"
					query="?${rest#*\?}"
					;;
				esac
				RESOLVED_DATABASE_URL="postgresql://${DATABASE_USER}:${DATABASE_PASSWORD}@${path}${query}"
				;;
			esac
			;;
		*)
			RESOLVED_DATABASE_URL="postgresql://${DATABASE_USER}:${DATABASE_PASSWORD}@${url}"
			;;
		esac

		case "$RESOLVED_DATABASE_URL" in
		*sslmode=*);;
		*rds.amazonaws.com*)
			case "$RESOLVED_DATABASE_URL" in
			*\?*) RESOLVED_DATABASE_URL="${RESOLVED_DATABASE_URL}&sslmode=require" ;;
			*) RESOLVED_DATABASE_URL="${RESOLVED_DATABASE_URL}?sslmode=require" ;;
			esac
			;;
		esac
	else
		RESOLVED_DATABASE_URL="$url"
	fi

	export RESOLVED_DATABASE_URL
}

_parse_database_url() {
	_resolve_database_url_env || return 1

	url="$RESOLVED_DATABASE_URL"
	case "$url" in
	postgresql://*) rest="${url#postgresql://}" ;;
	postgres://*) rest="${url#postgres://}" ;;
	*)
		echo "DATABASE_URL must use postgresql:// or postgres://" >&2
		return 1
		;;
	esac

	case "$rest" in
	*@*)
		userinfo="${rest%%@*}"
		hostpart="${rest#*@}"
		;;
	*)
		echo "could not parse database URL (missing credentials?)" >&2
		return 1
		;;
	esac

	DB_USER="${userinfo%%:*}"
	DB_PASSWORD="${userinfo#*:}"

	hostport="${hostpart%%/*}"
	pathquery="${hostpart#*/}"
	DB_NAME="${pathquery%%\?*}"

	case "$hostport" in
	*:*)
		DB_HOST="${hostport%%:*}"
		DB_PORT="${hostport#*:}"
		;;
	*)
		DB_HOST="$hostport"
		DB_PORT="5432"
		;;
	esac

	export DB_USER DB_PASSWORD DB_HOST DB_PORT DB_NAME
}

# Rewrite host/port for Docker networking (call after _parse_database_url).
database_url_rewrite_host() {
	new_host="$2"
	new_port="$3"

	query=""
	case "$RESOLVED_DATABASE_URL" in
	*\?*) query="?${RESOLVED_DATABASE_URL#*\?}" ;;
	esac

	printf 'postgresql://%s:%s@%s:%s/%s%s' "$DB_USER" "$DB_PASSWORD" "$new_host" "$new_port" "$DB_NAME" "$query"
}

_parse_database_url
