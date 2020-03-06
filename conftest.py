import os
import pytest
import sys

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), './libbeat/tests/system')))
import beat.compose


# Add support for Go-like tags
def pytest_addoption(parser):
    parser.addoption(
        "--tag",
        action="append",
        metavar="NAME",
        help="Enable tests marked with tag",
    )

def pytest_configure(config):
    config.addinivalue_line(
        "markers", "tag(name): Tag test with *name* using Go tag semantics"
    )


def pytest_runtest_setup(item):
    tags = [mark.args[0] for mark in item.iter_markers(name="tag")]
    if len(tags) == 0:
        return

    tag_opts = item.config.getoption("--tag") or [] 

    # Compatibility with environment variable
    add_opts_from_env(tag_opts,
        INTEGRATION_TESTS='integration',
    )
    if 'integration' in tag_opts:
        beat.compose.enable()

    if not tag_opts:
        pytest.skip()
        return

    for tag_opt in tag_opts:
        for tag in tags: 
            if tag not in tag_opts:
                pytest.skip()

def add_opts_from_env(opts, **kwargs):
    for variable, tag in kwargs.items():
        enabled = os.environ.get(variable, False)
        if enabled:
            opts.append(tag)
