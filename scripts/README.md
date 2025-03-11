# Script Management

## Legacy Scripts

This directory contains legacy scripts that have been consolidated into the new unified `kg.sh` script in the root directory.

The following scripts are now deprecated:
- `run.sh` - Use `./kg.sh start` instead
- `stop.sh` - Use `./kg.sh stop` instead
- `status.sh` - Use `./kg.sh status` instead

## New Unified Script

The new `kg.sh` script in the root directory provides a more organized and consistent interface for managing the Kay-Gee-Go application.

### Usage

```bash
# Start the application with default settings
./kg.sh start

# Start with custom settings
./kg.sh start --seed="Machine Learning" --max-nodes=200

# Stop the application
./kg.sh stop

# Restart the application
./kg.sh restart

# Show application status
./kg.sh status

# Run tests
./kg.sh test

# View logs
./kg.sh logs
./kg.sh logs --service=builder --follow

# Show help
./kg.sh help
```

For more details, run `./kg.sh help`. 