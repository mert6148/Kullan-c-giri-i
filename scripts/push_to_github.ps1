param(
    [string]$RemoteUrl = $env:GIT_REMOTE,
    [string]$Branch = $env:GIT_BRANCH,
    [string]$Message = "Update project files",
    [switch]$CreateRepo,
    [switch]$Force
)

Set-StrictMode -Version Latest

function Write-ErrExit($msg) {
    Write-Error $msg
    exit 1
}

if (-not (Get-Command git -ErrorAction SilentlyContinue)) {
    Write-ErrExit "Git not found. Please install git and try again."
}

if (-not $Branch) {
    # detect current branch
    $cur = git rev-parse --abbrev-ref HEAD 2>$null
    if ($LASTEXITCODE -eq 0) { $Branch = $cur.Trim() } else { $Branch = 'main' }
}

Write-Host "Using branch: $Branch"

if ($CreateRepo -and (Get-Command gh -ErrorAction SilentlyContinue)) {
    Write-Host "Creating repository via gh CLI (if not exists) and adding remote..."
    # If remote is a full URL, prefer that as target; otherwise attempt to create repo name
    if ($RemoteUrl) {
        # If remoteUrl is full https or ssh URL, still create remote via gh by converting
        git remote remove origin 2>$null | Out-Null
        git remote add origin $RemoteUrl
        git push -u origin $Branch
        exit $LASTEXITCODE
    } else {
        gh repo create --source=. --public --push --confirm
        $RemoteUrl = git remote get-url origin 2>$null
    }
}

if (-not $RemoteUrl) {
    # If remote missing, try to get origin or suggest user set it
    $exists = git remote get-url origin 2>$null
    if ($LASTEXITCODE -eq 0 -and $exists) {
        $RemoteUrl = $exists.Trim()
    } else {
        Write-Host "Remote origin not configured. Pass -RemoteUrl or set remote origin in git."
        Write-Host "To add a remote: git remote add origin <git@github.com:username/repo.git>"
        exit 1
    }
}

Write-Host "Using remote: $RemoteUrl"

# Stage all changes
git add -A
if ($LASTEXITCODE -ne 0) { Write-ErrExit "git add failed" }

# Check for changes
$status = git status --porcelain
if (-not $status) {
    Write-Host "No changes to commit."
    if (-not $Force) { Write-Host "Nothing to push"; exit 0 }
}

if ($status) {
    git commit -m $Message
    if ($LASTEXITCODE -ne 0) { Write-ErrExit "git commit failed" }
}

Write-Host "Pushing to origin/$Branch..."
git push -u origin $Branch
if ($LASTEXITCODE -ne 0) { Write-ErrExit "git push failed" }
Write-Host "Push completed."
