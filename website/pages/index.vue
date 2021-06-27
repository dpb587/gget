<template>
  <div class="page">
    <div class="pt-8 bg-teal-900 overflow-hidden md:py-16">
      <div class="max-w-xl mx-auto px-4 sm:px-6 md:px-8 md:max-w-screen-xl">
        <h1 class="text-center font-mono font-extrabold tracking-tight text-white sm:text-4xl text-3xl sm:leading-10 leading-8">
          gget
        </h1>
        <p class="mt-4 max-w-3xl mx-auto text-center text-xl leading-7 text-gray-300">
          a cli. to get things. from git repos.
        </p>
      </div>
      <div class="max-w-3xl mx-auto mt-8">
        <pre class="md:rounded border-b border-t md:border-l md:border-r border-teal-100"><code><span class="cmd-prompt">$ </span><strong>gget github.com/gohugoio/hugo 'hugo_extended_*_Linux-64bit.deb'</strong>
<div class="cmd-result">Found 1 file (13.9M) from github.com/gohugoio/hugo@v0.73.0
√ hugo_extended_0.73.0_Linux-64bit.deb done (sha256 OK)</div></code></pre>
      </div>
    </div>
    <div class="navbar py-2 md:py-1">
      <span class="ml-1 md:ml-2 text-gray-400">Jump to</span>
      <a href="#install" class="mx-1 md:mx-2"><fa icon="chevron-down" /> <span>Install</span></a>
      <a href="#introduction" class="mx-1 md:mx-2"><fa icon="play" /> <span>Getting Started</span></a>
      <a href="#advanced" class="mx-1 md:mx-2"><fa icon="graduation-cap" /> <span>Advanced</span></a>
      <a href="#cli" class="mx-1 md:mx-2"><fa icon="life-ring" /> <span>Reference</span></a>
      <a href="#docker" class="mx-1 md:mx-2"><fa :icon="['fab', 'docker']" /> <span>Docker</span></a>
    </div>
    <div class="max-w-4xl mx-auto">
      <div class="my-4 md:my-8 mx-4 md:mx-6">
        <ul class="md:grid md:grid-cols-2 md:col-gap-8 md:row-gap-6">
          <li v-for="fs in featureSpotlights" v-bind:key="fs.label" class="mt-4 md:mt-0">
            <div class="flex">
              <div class="flex-shrink-0">
                <div class="flex items-center justify-center h-12 w-12 rounded-md bg-teal-800 text-white">
                  <fa :icon="fs.icon" fixed-width class="h-6 w-6" />
                </div>
              </div>
              <div class="ml-4">
                <strong class="text-lg leading-6 font-medium text-gray-900" v-text="fs.label" />
                <p class="mt-1 text-base leading-5 text-gray-700" v-html="fs.description" />
              </div>
            </div>
          </li>
        </ul>
      </div>

      <div class="bg-white border border-gray-200 sm:rounded-lg sm:shadow mt-8 sm:mt-0">
        <div class="px-4 py-5 sm:px-6">
          <div class="-ml-4 -mt-4 flex justify-between items-center flex-wrap sm:flex-no-wrap">
            <div id="install" class="ml-4 pl-1 sm:pl-0 pt-3 sm:pt-4 pb-1 sm:pb-0">
              <h2 class="text-xl leading-7 font-medium text-gray-900">
                Installation
              </h2>
              <p class="mt-1 text-lg leading-6 text-gray-700">
                Version {{latest.origin.ref.replace(/^v/, '')}} (<a :href="`https://github.com/dpb587/gget/releases/tag/${latest.origin.ref}`">release notes</a>)
              </p>
            </div>
            <div class="ml-4 mt-2 md:mt-4 flex-shrink-0">
              <div class="relative inline-block text-left">
                <div @click="platformMenu = !platformMenu">
                  <span class="rounded-md shadow-sm">
                    <button type="button" class="inline-flex justify-center w-full rounded-md border border-teal-800 px-5 py-2 bg-white text-xl leading-8 font-medium text-teal-900 hover:text-teal-800 focus:outline-none focus:border-teal-600 focus:shadow-outline-teal active:bg-gray-50 active:text-teal-600 transition ease-in-out duration-150" aria-haspopup="true" aria-expanded="true">
                      <fa :icon="activePlatform.icon" class="-ml-1 mr-2 h-8 w-8" />
                      {{ activePlatform.label }}
                      <fa icon="chevron-down" class="ml-2 -mr-2 h-8 w-8" />
                    </button>
                  </span>
                </div>

                <div v-show="platformMenu" class="origin-top-right absolute right-0 mt-2 w-56 rounded-md shadow-lg z-10">
                  <div class="rounded-md bg-white shadow-xs" role="menu" aria-orientation="vertical" aria-labelledby="options-menu">
                    <div class="py-1">
                      <a v-for="(platform, name) in platformList" v-bind:key="name" :href="`#install-${name}`" @click.prevent="setPlatform(name)" class="group flex items-center px-4 py-2 text-lg leading-7 text-gray-700 hover:bg-gray-100 hover:text-gray-900 focus:outline-none focus:bg-gray-100 focus:text-gray-900" role="menuitem">
                        <fa :icon="platform.icon" fixed-width class="-ml-1 mr-2" />
                        {{ platform.label }}
                      </a>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div v-if="isUnixPlatform" class="px-4 py-5 border-t border-gray-300">
          <div class="flex flex-wrap sm:flex-no-wrap ">
            <div class="flex-shrink-0 pt-1">
              <fa icon="beer" fixed-width class="text-4xl text-teal-900" />
            </div>
            <div class="flex-grow ml-4">
              <div>
                <h3 class="text-lg leading-6 font-medium text-gray-900">
                  Homebrew
                </h3>
                <p class="mt-1 text-sm leading-5 text-gray-600">
                  Installing <code>gget</code> for your system.
                </p>
              </div>
            </div>
          </div>
          <div class="sm:-ml-1 sm:pl-16 mt-4">
            <pre><code><span class="cmd-prompt">$ </span>brew install dpb587/tap/gget</code></pre>
          </div>
        </div>
        <div v-if="isUnixPlatform" class="px-4 py-5 border-t border-gray-300">
          <div class="flex flex-wrap sm:flex-no-wrap">
            <div class="flex-shrink-0 pt-1">
              <fa icon="terminal" fixed-width class="text-4xl text-teal-900" />
            </div>
            <div class="flex-grow ml-4">
              <div>
                <h3 class="text-lg leading-6 font-medium text-gray-900">
                  Terminal
                </h3>
                <p class="mt-1 text-sm leading-5 text-gray-600">
                  Downloading to <code>./gget</code>.
                </p>
              </div>
            </div>
          </div>
          <div class="sm:-ml-1 sm:pl-16 mt-4">
            <pre><code><span class="cmd-prompt">$ </span>curl -Lo gget https://github.com/dpb587/gget/releases/download/{{latest.origin.ref}}/{{activePlatformResource.name}} \
  &amp;&amp; echo "<span v-text="activePlatformResource.checksums[0].data" />  gget" | shasum -c \
  &amp;&amp; chmod +x gget</code></pre>
          </div>
        </div>
        <div v-if="isWindowsPlatform" class="px-4 py-5 border-t border-gray-300">
          <div class="flex flex-wrap sm:flex-no-wrap">
            <div class="flex-shrink-0 pt-1">
              <fa icon="terminal" fixed-width class="text-4xl text-teal-900" />
            </div>
            <div class="flex-grow ml-4">
              <div>
                <h3 class="text-lg leading-6 font-medium text-gray-900">
                  PowerShell
                </h3>
                <p class="mt-1 text-sm leading-5 text-gray-600">
                  Downloading to <code>.\gget.exe</code>.
                </p>
              </div>
            </div>
          </div>
          <div class="sm:-ml-1 sm:pl-16 mt-4">
            <pre><code><span class="cmd-prompt">$ </span>( New-Object System.Net.WebClient ).DownloadFile("https://github.com/dpb587/gget/releases/download/{{latest.origin.ref}}/{{activePlatformResource.name}}", "$PWD\gget.exe")
<span class="cmd-prompt">$ </span>( Get-FileHash .\gget.exe -Algorithm {{ activePlatformResource.checksums[0].algo.toUpperCase() }} ).Hash -eq "{{ activePlatformResource.checksums[0].data.toUpperCase() }}"</code></pre>
          </div>
        </div>
        <div v-if="activePlatformName != 'other'" class="px-4 py-5 border-t border-gray-300">
          <div class="flex flex-wrap sm:flex-no-wrap">
            <div class="flex-shrink-0 pt-1">
              <fa icon="arrow-alt-circle-down" fixed-width class="text-4xl text-teal-900" />
            </div>
            <div class="flex-grow ml-4">
              <div class="flex-grow">
                <h3 class="text-lg leading-6 font-medium text-gray-900">
                  Download
                </h3>
                <p class="mt-1 text-sm leading-5 text-gray-600">
                  Verify (<span v-text="activePlatformResource.checksums[0].algo" />) and install yourself.
                </p>
              </div>
            </div>
            <div class="pl-16 -ml-1 mt-2 sm:pl-0 sm:ml-4 sm:mt-0 text-right">
              <a class="inline-flex rounded-md shadow-sm" :href="`https://github.com/dpb587/gget/releases/download/${latest.origin.ref}/${activePlatformResource.name}`">
                <button type="button" class="inline-flex items-center px-4 py-2 border border-transparent text-base leading-6 font-medium rounded-md text-white bg-teal-900 hover:bg-teal-800 focus:outline-none focus:border-teal-700 focus:shadow-outline-teal active:bg-teal-800 transition ease-in-out duration-150">
                  <fa icon="file-download" class="mr-3" />
                  <span v-text="activePlatformResource.name" />
                </button>
              </a>
            </div>
          </div>
          <div class="sm:-ml-1 sm:pl-16 mt-4">
            <pre><code><span v-text="activePlatformResource.checksums[0].data" />  <span v-text="activePlatformResource.name" /></code></pre>
          </div>
        </div>
        <div v-if="activePlatformName == 'other'" class="px-4 py-5 border-t border-gray-300">
          <div class="flex flex-wrap sm:flex-no-wrap">
            <div class="flex-shrink-0 pt-1">
              <fa icon="tools" fixed-width class="text-4xl text-teal-900" />
            </div>
            <div class="flex-grow ml-4">
              <div>
                <h3 class="text-lg leading-6 font-medium text-gray-900">
                  Build with <a href="https://golang.org/">go</a>
                </h3>
                <p class="mt-1 text-sm leading-5 text-gray-600">
                  Compile into a local executable.
                </p>
              </div>
            </div>
          </div>
          <div class="sm:-ml-1 sm:pl-16 mt-4">
            <pre><code><span class="cmd-prompt">$ </span>git clone https://github.com/dpb587/gget.git
<span class="cmd-prompt">$ </span>cd gget/
<span class="cmd-prompt">$ </span>git checkout {{ latest.origin.ref }}
<span class="cmd-prompt">$ </span>go build .</code></pre>
          </div>
        </div>
      </div>
      
      <div>
        <div id="introduction" class="text-center mt-4 sm:mt-6 pt-8 sm:pt-12">
          <h2 class="text-3xl leading-8 tracking-tight font-bold text-gray-900 sm:text-4xl sm:leading-6">
            Getting Started
          </h2>
          <p class="mt-1 max-w-2xl mx-auto text-xl leading-6 text-gray-700 sm:mt-4">
            An introduction of the <code>gget</code> concepts.
          </p>
        </div>

        <div class="px-4 mt-12 md:mt-16">
          <h2 class="text-2xl leading-7 font-medium text-black">
            Pass the Origin
          </h2>
          <p class="mt-3 text-lg leading-7 text-gray-800">
            The only required argument to <code>gget</code> is the download origin &ndash; it is a shortened URL specifying where to find the repository. By default, the latest version will be discovered and all release assets will be downloaded.
          </p>
          <pre class="my-3 -mx-4"><code><span class="cmd-prompt">$ </span>gget <strong>github.com/dpb587/gget</strong>
<div class="cmd-result">Found 3 files (33M) from github.com/dpb587/gget@{{latest.origin.ref}}
√ gget-{{latestOriginRefSemver}}-darwin-amd64      done (sha256 OK)
√ gget-{{latestOriginRefSemver}}-linux-amd64       done (sha256 OK)
√ gget-{{latestOriginRefSemver}}-windows-amd64.exe done (sha256 OK)</div></code></pre>
          <p class="mt-3 text-lg leading-7 text-gray-800">
            If you want something other than latest, add an <strong><code>@REF</code></strong> to the end. For example, to always download from the <code v-text="latest.origin.ref" /> release tag you would use <code>github.com/dpb587/gget@{{latest.origin.ref}}</code>.
          </p>
        </div>

        <div class="px-4 mt-12 md:mt-16">
          <h2 class="text-2xl leading-7 font-medium text-black">
            Name the Resources
          </h2>
          <p class="mt-3 text-lg leading-7 text-gray-800">
            To avoid downloading all resources, list only the names you want as additional arguments. Basic glob patterns like wildcards (<code>*</code>) can be used if you don't know the exact name.
          </p>
          <pre class="my-3 -mx-4"><code><span class="cmd-prompt">$ </span>gget github.com/dpb587/gget <strong>'*linux*' gget-{{latestOriginRefSemver}}-darwin-amd64</strong>
<div class="cmd-result">Found 2 files (22.4M) from github.com/dpb587/gget@{{latest.origin.ref}}
√ gget-{{latestOriginRefSemver}}-darwin-amd64 done (sha256 OK)
√ gget-{{latestOriginRefSemver}}-linux-amd64  done (sha256 OK)</div></code></pre>
          <p class="mt-3 text-lg leading-7 text-gray-800">
            If patterns may match more resources than desired, use the <strong><code>--exclude</code></strong> option to restrict further.
          </p>
          <pre class="my-3 -mx-4"><code><span class="cmd-prompt">$ </span>gget github.com/dpb587/gget <strong>--exclude='*-amd64'</strong>
<div class="cmd-result">Found 1 file (10.6M) from github.com/dpb587/gget@{{latest.origin.ref}}
√ gget-{{latestOriginRefSemver}}-windows-amd64.exe done (sha256 OK)</div></code></pre>
          <p class="mt-3 text-lg leading-7 text-gray-800">
            If an argument does not end up matching anything, the command will error before starting any downloads. If missing resources are expected, the <strong><code>--ignore-missing</code></strong> option may be used to ignore that error.
          </p>
          <pre class="my-3 -mx-4"><code><span class="cmd-prompt">$ </span>gget github.com/dpb587/gget <strong>--ignore-missing</strong> '*-darwin-*' '*-macos-*'
<div class="cmd-result">Found 1 file (11.7M) from github.com/dpb587/gget@{{latest.origin.ref}}
√ gget-{{latestOriginRefSemver}}-darwin-amd64 done (sha256 OK)</div></code></pre>
        </div>

        <div class="px-4 mt-12 md:mt-16">
          <h2 class="text-2xl leading-7 font-medium text-black">
            Remote Servers
          </h2>
          <p class="mt-3 text-lg leading-7 text-gray-800">
            These examples have been using public repositories from <a href="https://github.com/">GitHub</a>, but <code>gget</code> also supports <a href="https://gitlab.com/">GitLab</a> and self-hosted installations of the two. If the service is not automatically detected, the <strong><code>--service</code></strong> option may be used to override the API used to either <code>github</code> or <code>gitlab</code>.
          </p>
          <p class="mt-3 text-lg leading-7 text-gray-800">
            To access private repositories, be sure to configure your access token through the appropriate environment variable or a machine password in your <a href="https://ec.haxx.se/usingcurl/usingcurl-netrc">~/.netrc</a> file. For GitHub, the <strong><code>$GITHUB_TOKEN</code></strong> environment variable is used, and for GitLab it is <strong><code>$GITLAB_TOKEN</code></strong>.
          </p>
          <pre class="my-3 -mx-4"><code><span class="cmd-prompt">$ </span><strong>export GITLAB_TOKEN=</strong>a1b2c3d4...
<span class="cmd-prompt">$ </span>gget <strong>--service=gitlab</strong> gitlab.acme.corp/my-team/product '*-linux-amd64'
<div class="cmd-result">Found 1 file from gitlab.acme.corp/my-team/product@v8.21.2
√ product-8.21.2-darwin-amd64 done</div></code></pre>
        </div>

        <div class="px-4 mt-12 md:mt-16">
          <h2 class="text-2xl leading-7 font-medium text-black">
            Blobs and Archives
          </h2>
          <p class="mt-3 text-lg leading-7 text-gray-800">
            In addition to user-uploaded release assets, <code>gget</code> also supports downloading individual source files from a repository. The <strong><code>--type</code></strong> option can be used to switch to <code>blob</code>-mode.
          </p>
          <pre class="my-3 -mx-4"><code><span class="cmd-prompt">$ </span>gget github.com/dpb587/gget <strong>--type=blob README.md</strong>
<div class="cmd-result">Found 1 file (6.8K) from github.com/dpb587/gget@{{latest.origin.ref}}
√ README.md done</div></code></pre>
          <p class="mt-3 text-lg leading-7 text-gray-800">
            Alternatively, to get an archive of the entire repository, use the <code>archive</code> type.
          </p>
          <pre class="my-3 -mx-4"><code><span class="cmd-prompt">$ </span>gget github.com/dpb587/gget <strong>--type=archive '*.zip'</strong>
<div class="cmd-result">Found 1 file from github.com/dpb587/gget@{{latest.origin.ref}}
√ gget-{{latest.origin.ref}}.zip done</div></code></pre>
          <p class="mt-3 text-lg leading-7 text-gray-800">
            For these source types, branches and commits may also be used as a reference.
          </p>
          <pre class="my-3 -mx-4"><code><span class="cmd-prompt">$ </span>gget github.com/dpb587/gget<strong>@d75361c</strong> --type=archive '*.zip'
<div class="cmd-result">Found 1 file from github.com/dpb587/gget@d75361c1316c6004cde13c04460c57552a839d57
√ gget-d75361c13.zip done</div></code></pre>
        </div>

        <div class="px-4 mt-12 md:mt-16">
          <h2 class="text-2xl leading-7 font-medium text-black">
            Troubleshooting
          </h2>
          <p class="mt-3 text-lg leading-7 text-gray-800">
            If you are unable to find the resources you expect, try verifying <code>gget</code> sees them with the <strong><code>--list</code></strong> option.
          </p>
          <pre class="my-3 -mx-4"><code><span class="cmd-prompt">$ </span>gget github.com/dpb587/gget --type=archive <strong>--list</strong>
<div class="cmd-result">Found 2 files from github.com/dpb587/gget@{{latest.origin.ref}}
gget-{{latest.origin.ref}}.tar.gz
gget-{{latest.origin.ref}}.zip</div></code></pre>
          <p class="mt-3 text-lg leading-7 text-gray-800">
            If <code>gget</code> is erroring or not behaving as you expect, try enabling logging with the <strong><code>-v</code></strong> verbosity flag. If logs are unhelpful or you still have a concern, refer to them when searching and reporting a <a href="https://github.com/dpb587/gget/issues">project issue</a>.
          </p>
          <pre class="my-3 -mx-4"><code><span class="cmd-prompt">$ </span>gget github.com/dpb587/gget <strong>-vvv</strong> --no-progress '*linux*'</code></pre>
        </div>
      </div>
      
      <div>
        <div id="advanced" class="text-center mt-4 sm:mt-6 pt-8 sm:pt-12">
          <h2 class="text-3xl leading-8 tracking-tight font-bold text-gray-900 sm:text-4xl sm:leading-6">
            Advanced
          </h2>
          <p class="max-w-2xl mx-auto text-xl leading-6 text-gray-700 sm:mt-4">
            More behaviors and tricks of <code>gget</code>.
          </p>
        </div>

        <div class="px-4 mt-12 md:mt-16">
          <h2 class="text-2xl leading-7 font-medium text-black">
            Finding Versions
          </h2>
          <p class="mt-3 text-lg leading-7 text-gray-800">
            By default, releases tagged as "pre-release" are not used when finding the latest version. To find the latest pre-release, use the <strong><code>--ref-stability</code></strong> option and specify <code>pre-release</code> (or, if you don't care whether it is pre-release or stable, use <code>any</code>).
          </p>
          <pre class="my-3 -mx-4"><code><span class="cmd-prompt">$ </span>gget github.com/dpb587/gget-test <strong>--ref-stability=pre-release</strong> '*linux*'
<div class="cmd-result">Found 1 file (11.7M) from github.com/dpb587/gget-test@v0.5.0-rc.1
√ gget-0.5.0-rc.1-linux-amd64 done (sha256 OK)</div></code></pre>
          <p class="mt-3 text-lg leading-7 text-gray-800">
            Sometimes, instead of a literal version, you may want to use a semver constraint to continue following the latest patches of a version. The <strong><code>--ref-version</code></strong> option allows you to specify a set of <a href="https://github.com/Masterminds/semver#checking-version-constraints">version constraints</a>. The first semver-parseable tag which matches will be used.
          </p>
          <pre class="my-3 -mx-4"><code><span class="cmd-prompt">$ </span>gget github.com/dpb587/gget <strong>--ref-version={{latestOriginRefSemver.replace(/(\.[^\.]+)$/, '.x')}}</strong> '*linux*'
<div class="cmd-result">Found 1 file (11.7M) from github.com/dpb587/gget@{{latest.origin.ref}}
√ gget-{{latestOriginRefSemver}}-linux-amd64 done (sha256 OK)</div></code></pre>
        </div>

        <div class="px-4 mt-12 md:mt-16">
          <h2 class="text-2xl leading-7 font-medium text-black">
            Manipulating Files
          </h2>
          <p class="mt-3 text-lg leading-7 text-gray-800">
            To save a resource under a different local name (or path), prefix it to the resource name with an equals (<strong><code>=</code></strong>).
          </p>
          <pre class="my-3 -mx-4"><code><span class="cmd-prompt">$ </span>gget github.com/dpb587/gget <strong>/tmp/gget=</strong>'*linux*'
<div class="cmd-result">Found 1 file (11.7M) from github.com/dpb587/gget@{{latest.origin.ref}}
√ gget-0.4.0-linux-amd64 done (sha256 OK)</div></code></pre>
          <p class="mt-3 text-lg leading-7 text-gray-800">
            To automatically mark a resource (such as a binary) as executable, use the <strong><code>--executable</code></strong> option.
          </p>
          <pre class="my-3 -mx-4"><code><span class="cmd-prompt">$ </span>gget github.com/dpb587/gget <strong>--executable</strong> '*linux*'
<div class="cmd-result">Found 1 file (11.7M) from github.com/dpb587/gget@{{latest.origin.ref}}
√ gget-0.4.0-linux-amd64 done (sha256 OK; executable)</div></code></pre>
          <p class="mt-3 text-lg leading-7 text-gray-800">
            When working with archive resources, you may want to use the <strong><code>--stdout</code></strong> option. This causes all downloaded contents to be streamed to <code>STDOUT</code> for redirects or chained commands.
          </p>
          <pre class="my-3 -mx-4"><code><span class="cmd-prompt">$ </span>gget github.com/dpb587/gget --type=archive '*.tar.gz' <strong>--stdout | tar -tzf-</strong> | tail -n2
<div class="cmd-result">Found 1 file from github.com/dpb587/gget@{{latest.origin.ref}}
√ gget-{{latest.origin.ref}}.tar.gz done
dpb587-gget-{{latest.metadata.filter((v) => v.key == 'commit')[0].value.substr(0, 7)}}/scripts/integration.test.sh
dpb587-gget-{{latest.metadata.filter((v) => v.key == 'commit')[0].value.substr(0, 7)}}/scripts/website.build.sh</div></code></pre>
        </div>

        <div class="px-4 mt-12 md:mt-16">
          <h2 class="text-2xl leading-7 font-medium text-black">
            Verifying Checksums
          </h2>
          <p class="mt-3 text-lg leading-7 text-gray-800">
            By default, <code>gget</code> attempts to discover and verify checksums (based on conventional file names or release notes). If a checksum is not found, no verification is attempted; but if one is found and it fails, the command will error. Use the <strong><code>--verify-checksum</code></strong> option to require a specific algorithm, minimum algorithm strength, or disable the behavior entirely.
          </p>
          <pre class="my-3 -mx-4"><code><span class="cmd-prompt">$ </span>gget <strong>--verify-checksum=sha512-min</strong> github.com/dpb587/gget
<div class="cmd-result">Found 1 file (11.7M) from github.com/dpb587/gget@{{latest.origin.ref}} '*linux*'
gget: error: preparing transfer of gget-0.4.0-linux-amd64: acceptable checksum required but not found: sha512</div></code></pre>
        </div>

        <div class="px-4 mt-12 md:mt-16">
          <h2 class="text-2xl leading-7 font-medium text-black">
            Script Integrations
          </h2>
          <p class="mt-3 text-lg leading-7 text-gray-800">
            The only output sent to <code>STDOUT</code> is from downloads configured to write there (or when the <code>--export</code> option is used). All other information (status updates, progress, and log messages) will be written to <code>STDERR</code>. If you want to reduce runtime output or rely on your own messaging, use the <strong><code>--quiet</code></strong> option.
          </p>
          <p class="mt-3 text-lg leading-7 text-gray-800">
            The <strong><code>--export</code></strong> option is helpful for capturing details about matched resources, the resolved origin, and some additional metadata. It supports output for JSON (<a href="latest.json">example</a>), JSONPath-like (<a href="https://kubernetes.io/docs/reference/kubectl/jsonpath/">docs</a>), YAML, and a plain format.
          </p>
          <pre class="my-3 -mx-4"><code><span class="cmd-prompt">$ </span>gget github.com/dpb587/gget <strong>--export=jsonpath='{.origin.ref}'</strong> --no-download --quiet
<div class="cmd-result"><span v-text="latest.origin.ref" /></div></code></pre>
          <pre class="my-3 -mx-4"><code><span class="cmd-prompt">$ </span>gget github.com/dpb587/gget <strong>--export=json</strong> --no-download --quiet \
  | jq '.metadata | from_entries'
<div class="cmd-result">{
<div v-for="metadatum in latest.metadata.filter((v) => v.key != 'github-release-body')" v-bind:key="metadatum.key" v-text="`  ${JSON.stringify(metadatum.key)}: ${JSON.stringify(metadatum.value)},`" />}</div></code></pre>
          <pre class="my-3 -mx-4"><code><span class="cmd-prompt">$ </span>gget github.com/dpb587/gget <strong>--export=plain</strong> --no-download --quiet \
  | awk '{ if ( $1 == "resource-checksum" ) { print $4 "  " $2 } }'
<div class="cmd-result"><div v-for="resource in latest.resources" v-bind:key="resource.name" v-text="`${resource.checksums[0].data}  ${resource.name}`" /></div></code></pre>
          <p class="mt-3 text-lg leading-7 text-gray-800">
            The <strong><code>--version</code></strong> option can optionally be passed a constraint if you need to rely on a particular version of <code>gget</code> itself. When the constraint is not met, the command will exit with an error.
          </p>
          <pre class="my-3 -mx-4"><code>if ! gget --version<strong>='>=0.5'</strong> >/dev/null 2>&amp;1
then
  echo "oops. gget version 0.5 or greater must be installed first." >&amp;2
  exit 1
fi</code></pre>
        </div>
      </div>
      
      <div>
        <div id="cli" class="text-center mt-4 sm:mt-6 pt-8 sm:pt-12">
          <h2 class="text-3xl leading-8 tracking-tight font-bold text-gray-900 sm:text-4xl sm:leading-6">
            Reference
          </h2>
          <p class="max-w-2xl mx-auto text-xl leading-6 text-gray-700 sm:mt-4">
            Showing <code>gget --help</code> from <code v-text="latest.origin.ref" />.
          </p>
        </div>
        <div class="px-4 mt-8">
          <pre class="my-3 -mx-4"><code v-text="helptext.replace(/\s+$/g, '')" /></pre>
        </div>
      </div>

      <div>
        <div id="docker" class="text-center mt-4 sm:mt-6 pt-8 sm:pt-12">
          <h2 class="text-3xl leading-8 tracking-tight font-bold text-gray-900 sm:text-4xl sm:leading-6">
            Docker
          </h2>
          <p class="max-w-2xl mx-auto text-xl leading-6 text-gray-700 sm:mt-4">
            Using <code>gget</code> to build containers.
          </p>
        </div>
        <div class="px-4 mt-8">
          <p class="mt-3 text-lg leading-7 text-gray-800">
            A <code>gget</code> Docker image is available and easily integrates as a build stage for downloads. It includes tools for uncompressing archives, and the default working directory is <code>/result</code> for access in later stages.
          </p>
          <pre class="my-3 -mx-4"><code>FROM dpb587/gget as gget
RUN gget github.com/cloudfoundry/bosh-cli --ref-version=5.x \
      --executable bosh=bosh-cli-*-linux-amd64
RUN gget github.com/cloudfoundry/bosh-bootloader \
      --executable bbl=bbl-*_linux_x86-64
RUN gget github.com/pivotal-cf/om \
      --stdout om-linux-*.tar.gz | tar -xzf- om

FROM ubuntu
COPY --from=gget /result/* /usr/local/bin/
# ...everything else for your image...</code></pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from '~/plugins/axios'

export default {
  head: {
    title: 'gget // a cli. to get things. from git repos.',
    meta: [
      { hid: 'description', name: 'description', content: 'An open-source tool that makes it easier to find and automate file downloads from GitHub and GitLab repositories.' },
      { hid: 'twitter:card', name: 'twitter:card', content: 'summary_large_image' },
      { hid: 'twitter:creator', name: 'twitter:creator', content: '@dpb587' },
      { hid: 'twitter:title', name: 'twitter:title', content: 'gget // a cli. to get things. from git repos.' },
      { hid: 'twitter:description', name: 'twitter:description', content: 'An open-source tool that makes it easier to find and automate file downloads from GitHub and GitLab repositories.' },
      { hid: 'twitter:image', name: 'twitter:image', content: `${process.env.baseUrl}img/card~v1~1280x640.png` },
      { hid: 'twitter:image:alt', name: 'twitter:image:alt', content: 'The gget name, tagline, and icons representing key features.'}
    ]
  },
  async asyncData() {
    const { data: latest } = await axios.get('latest.json')
    const { data: helptext } = await axios.get('latest-help.txt')

    return {
      latest: latest,
      helptext: helptext,
    }
  },
  mounted() {
    // naive checks
    if (navigator.appVersion.indexOf('Linux') > -1) {
      this.activePlatformName = 'linux-amd64'
    } else if (navigator.appVersion.indexOf('Mac') > -1) {
      this.activePlatformName = 'macos-amd64'
    } else if (navigator.appVersion.indexOf('Win') > -1) {
      this.activePlatformName = 'windows-amd64'
    }
  },
  data() {
    return {
      activePlatformName: 'other',
      platformMenu: false,
      platformList: {
        'linux-amd64': {
          label: 'Linux (amd64)',
          icon: ['fab', 'linux'],
          resourceSuffix: 'linux-amd64',
        },
        'macos-amd64': {
          label: 'macOS (amd64)',
          icon: ['fab', 'apple'],
          resourceSuffix: 'darwin-amd64',
        },
        'windows-amd64': {
          label: 'Windows (amd64)',
          icon: ['fab', 'windows'],
          resourceSuffix: 'windows-amd64.exe',
        },
        'other': {
          label: 'Other (source)',
          icon: 'code'
        },
      },
      featureSpotlights: [
        {
          icon: 'play-circle',
          label: 'Standalone CLI',
          description: 'No <code>git</code> or local clones required.'
        },
        {
          icon: 'server',
          label: 'Public & Private Servers',
          description: 'Supporting both <a href="https://github.com/">GitHub</a> and <a href="https://gitlab.com/">GitLab</a>.'
        },
        {
          icon: 'lock',
          label: 'Public & Private Repos',
          description: 'API-based access for private resources.'
        },
        {
          icon: 'code-branch',
          label: 'Tags, Branches, and Commits',
          description: 'Also semver-based constraint matching.'
        },
        {
          icon: 'copy',
          label: 'Archives, Assets, and Blobs',
          description: 'Download any type of resource from repos.'
        },
        {
          icon: 'random',
          label: 'Built-in File Operations',
          description: 'Easy rename, binary, and stream options.'
        },
        {
          icon: 'check-double',
          label: 'Checksum Verification',
          description: 'Verify downloads by SHA (if published).'
        },
        {
          icon: 'heart',
          label: 'Open-Source Project',
          description: 'Contribute and discuss at <a href="https://github.com/dpb587/gget">dpb587/gget</a>.'
        }
      ]
    }
  },
  computed: {
    activePlatform() {
      return this.platformList[this.activePlatformName]
    },
    activePlatformResource() {
      return this.latest.resources.filter((v) => v.name.indexOf(this.activePlatform.resourceSuffix) > 0)[0]
    },
    isWindowsPlatform() {
      switch (this.activePlatformName) {
        case 'windows-amd64':
          return true
      }

      return false
    },
    isUnixPlatform() {
      switch (this.activePlatformName) {
        case 'linux-amd64':
        case 'macos-amd64':
          return true
      }

      return false
    },
    latestOriginRefSemver() {
      return this.latest.origin.ref.replace(/^v/, '');
    }
  },
  methods: {
    setPlatform(name) {
      this.activePlatformName = name;
      this.platformMenu = false
    }
  }
}
</script>

<style>
.navbar {
  @apply px-2 bg-teal-800 text-center text-gray-200 leading-8;
}

.navbar > * {
  @apply px-2 py-1;
}

.navbar > a {
  @apply whitespace-no-wrap;
}

.navbar > a:hover {
  @apply text-white;
}

.navbar > a > span {
  @apply ml-1 underline;
}

pre {
  @apply bg-teal-900 px-4 py-3 text-white overflow-auto;
}

pre code strong {
  @apply font-bold;
}

pre code .cmd-result {
  @apply text-gray-400;
}

pre code .cmd-prompt {
  @apply text-teal-400;
  user-select: none;
}

p a {
  @apply underline text-teal-900;
}

p code {
  @apply text-teal-900;
}

pre + p {
  @apply mt-6 !important;
}
</style>

<!-- the nicest lil 'ol web page an unknown cli didn't know it needed -->
