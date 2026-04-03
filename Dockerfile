FROM node:22-slim AS base
RUN corepack enable pnpm

# Install dependencies
FROM base AS deps
WORKDIR /app
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY app/package.json app/
COPY packages/core/package.json packages/core/
RUN pnpm install --frozen-lockfile

# Build
FROM base AS build
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY --from=deps /app/app/node_modules ./app/node_modules
COPY --from=deps /app/packages/core/node_modules ./packages/core/node_modules
COPY . .
RUN pnpm --filter app run build

# Production
FROM base AS production
WORKDIR /app
ENV NODE_ENV=production
ENV PORT=3000
ENV PROMETHEUS_VAULT_PATH=/data/vault
ENV PROMETHEUS_DB_PATH=/data/prometheus.db

# Copy built app
COPY --from=build /app/app/build ./build
COPY --from=build /app/app/package.json ./
COPY --from=deps /app/app/node_modules ./node_modules

# Create data directory
RUN mkdir -p /data/vault/.prometheus

EXPOSE 3000

CMD ["node", "build/index.js"]
